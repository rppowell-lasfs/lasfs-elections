from flask import Flask, request, jsonify, session, redirect, url_for
from flask_sqlalchemy import SQLAlchemy
from werkzeug.security import generate_password_hash, check_password_hash
from functools import wraps

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///:memory:'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
app.config['SECRET_KEY'] = 'supersecretkey'
db = SQLAlchemy(app)

# Models
class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    password_hash = db.Column(db.String(128), nullable=False)
    votes = db.relationship('Vote', backref='user', lazy=True)

    def set_password(self, password):
        self.password_hash = generate_password_hash(password)

    def check_password(self, password):
        return check_password_hash(self.password_hash, password)

class Election(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(120), nullable=False)
    candidates = db.relationship('Candidate', backref='election', lazy=True)
    votes = db.relationship('Vote', backref='election', lazy=True)

class Candidate(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(120), nullable=False)
    election_id = db.Column(db.Integer, db.ForeignKey('election.id'), nullable=False)

class Vote(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('user.id'), nullable=False)
    election_id = db.Column(db.Integer, db.ForeignKey('election.id'), nullable=False)
    # Store ranked candidate ids as comma separated string: "3,1,2"
    ranked_choices = db.Column(db.String, nullable=False)

# Helpers
def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if 'user_id' not in session:
            return jsonify({"error": "Authentication required"}), 401
        return f(*args, **kwargs)
    return decorated_function

def runoff_winner(election):
    # Collect votes as lists of candidate ids in ranked order
    votes = []
    candidate_ids = [c.id for c in election.candidates]
    for vote in election.votes:
        ranked = list(map(int, vote.ranked_choices.split(',')))
        votes.append(ranked)

    # Runoff voting algorithm
    remaining = set(candidate_ids)
    majority = len(votes) / 2

    while True:
        # Count first-choice votes for remaining candidates
        counts = {cid: 0 for cid in remaining}
        for vote in votes:
            for choice in vote:
                if choice in remaining:
                    counts[choice] += 1
                    break

        # Check if any candidate has majority
        for cid, count in counts.items():
            if count > majority:
                return Candidate.query.get(cid).name

        # Find candidate(s) with fewest votes
        min_votes = min(counts.values())
        to_eliminate = [cid for cid, count in counts.items() if count == min_votes]

        # If tie among all remaining candidates, return all as winners
        if len(to_eliminate) == len(remaining):
            return [Candidate.query.get(cid).name for cid in remaining]

        # Eliminate candidate(s) with fewest votes
        for cid in to_eliminate:
            remaining.remove(cid)

# Routes
@app.route('/register', methods=['POST'])
def register():
    data = request.json
    if not data or not data.get('username') or not data.get('password'):
        return jsonify({"error": "Username and password required"}), 400
    if User.query.filter_by(username=data['username']).first():
        return jsonify({"error": "Username already exists"}), 400
    user = User(username=data['username'])
    user.set_password(data['password'])
    db.session.add(user)
    db.session.commit()
    return jsonify({"message": "User registered successfully"})

@app.route('/login', methods=['POST'])
def login():
    data = request.json
    if not data or not data.get('username') or not data.get('password'):
        return jsonify({"error": "Username and password required"}), 400
    user = User.query.filter_by(username=data['username']).first()
    if user and user.check_password(data['password']):
        session['user_id'] = user.id
        return jsonify({"message": "Logged in successfully"})
    return jsonify({"error": "Invalid credentials"}), 401

@app.route('/logout', methods=['POST'])
@login_required
def logout():
    session.pop('user_id', None)
    return jsonify({"message": "Logged out successfully"})

@app.route('/elections', methods=['POST'])
@login_required
def create_election():
    data = request.json
    if not data or not data.get('name') or not data.get('candidates'):
        return jsonify({"error": "Election name and candidates required"}), 400
    election = Election(name=data['name'])
    db.session.add(election)
    db.session.flush()
    for cname in data['candidates']:
        candidate = Candidate(name=cname, election_id=election.id)
        db.session.add(candidate)
    db.session.commit()
    return jsonify({"message": "Election created", "election_id": election.id})

@app.route('/elections/<int:election_id>/vote', methods=['POST'])
@login_required
def vote(election_id):
    data = request.json
    user_id = session['user_id']
    election = Election.query.get_or_404(election_id)
    if not data or not data.get('ranked_choices'):
        return jsonify({"error": "Ranked choices required"}), 400

    ranked_choices = data['ranked_choices']
    if not isinstance(ranked_choices, list) or not ranked_choices:
        return jsonify({"error": "Ranked choices must be a non-empty list"}), 400

    candidate_ids = {c.id for c in election.candidates}
    if set(ranked_choices) - candidate_ids:
        return jsonify({"error": "Invalid candidate IDs in ranked choices"}), 400

    # Check if user already voted in this election
    existing_vote = Vote.query.filter_by(user_id=user_id, election_id=election_id).first()
    if existing_vote:
        return jsonify({"error": "User has already voted in this election"}), 400

    vote = Vote(user_id=user_id, election_id=election_id, ranked_choices=','.join(map(str, ranked_choices)))
    db.session.add(vote)
    db.session.commit()
    return jsonify({"message": "Vote recorded"})

@app.route('/elections/<int:election_id>/result', methods=['GET'])
@login_required
def result(election_id):
    election = Election.query.get_or_404(election_id)
    winner = runoff_winner(election)
    return jsonify({"election": election.name, "winner": winner})

@app.route('/elections', methods=['GET'])
@login_required
def list_elections():
    elections = Election.query.all()
    result = []
    for e in elections:
        result.append({
            "id": e.id,
            "name": e.name,
            "candidates": [c.name for c in e.candidates]
        })
    return jsonify(result)

if __name__ == '__main__':
    with app.app_context():
        db.create_all()
    app.run(debug=True)
