#
from flask import Flask, render_template_string, request, redirect, url_for, session, flash
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
    # Store ranked candidate ids as comma separated string
    ranked_choices = db.Column(db.String, nullable=False)

# Helpers
def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if 'user_id' not in session:
            flash("Login required")
            return redirect(url_for('login'))
        return f(*args, **kwargs)
    return decorated_function

def runoff_winner(election):
    # Collect all votes for election
    votes = [vote.ranked_choices.split(',') for vote in election.votes]
    candidates = [str(c.id) for c in election.candidates]
    total_votes = len(votes)
    if total_votes == 0:
        return None

    # Runoff rounds
    while True:
        # Count first-choice votes
        counts = {cid: 0 for cid in candidates}
        for vote in votes:
            for choice in vote:
                if choice in candidates:
                    counts[choice] += 1
                    break

        # Check if any candidate has majority
        for cid, count in counts.items():
            if count > total_votes / 2:
                return Candidate.query.get(int(cid))

        # Find candidate(s) with fewest votes
        min_votes = min(counts.values())
        to_remove = [cid for cid, count in counts.items() if count == min_votes]

        # Remove candidate(s) with fewest votes
        candidates = [cid for cid in candidates if cid not in to_remove]

        # If tie or no candidates left
        if len(candidates) == 0:
            return None
        if len(candidates) == 1:
            return Candidate.query.get(int(candidates[0]))

        # Remove eliminated candidates from votes
        new_votes = []
        for vote in votes:
            new_vote = [choice for choice in vote if choice in candidates]
            new_votes.append(new_vote)
        votes = new_votes

# Routes and views
@app.route('/')
def home():
    elections = Election.query.all()
    user = None
    if 'user_id' in session:
        user = User.query.get(session['user_id'])
    return render_template_string('''
    <!doctype html>
    <title>Runoff Voting Elections</title>
    <h1>Welcome{% if user %}, {{ user.username }}{% endif %}</h1>
    {% if user %}
        <a href="{{ url_for('logout') }}">Logout</a>
    {% else %}
        <a href="{{ url_for('login') }}">Login</a> | <a href="{{ url_for('register') }}">Register</a>
    {% endif %}
    <h2>Elections</h2>
    <ul>
    {% for election in elections %}
        <li><a href="{{ url_for('election_detail', election_id=election.id) }}">{{ election.name }}</a></li>
    {% else %}
        <li>No elections available.</li>
    {% endfor %}
    </ul>
    ''', elections=elections, user=user)

@app.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        username = request.form['username'].strip()
        password = request.form['password']
        if not username or not password:
            flash("Username and password required")
            return redirect(url_for('register'))
        if User.query.filter_by(username=username).first():
            flash("Username already taken")
            return redirect(url_for('register'))
        user = User(username=username)
        user.set_password(password)
        db.session.add(user)
        db.session.commit()
        flash("Registration successful. Please login.")
        return redirect(url_for('login'))
    return render_template_string('''
    <!doctype html>
    <title>Register</title>
    <h1>Register</h1>
    <form method="post">
      Username: <input type="text" name="username" required><br>
      Password: <input type="password" name="password" required><br>
      <input type="submit" value="Register">
    </form>
    <a href="{{ url_for('login') }}">Login</a>
    ''')

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        username = request.form['username'].strip()
        password = request.form['password']
        user = User.query.filter_by(username=username).first()
        if user and user.check_password(password):
            session['user_id'] = user.id
            flash("Logged in successfully")
            return redirect(url_for('home'))
        flash("Invalid username or password")
        return redirect(url_for('login'))
    return render_template_string('''
    <!doctype html>
    <title>Login</title>
    <h1>Login</h1>
    <form method="post">
      Username: <input type="text" name="username" required><br>
      Password: <input type="password" name="password" required><br>
      <input type="submit" value="Login">
    </form>
    <a href="{{ url_for('register') }}">Register</a>
    ''')

@app.route('/logout')
def logout():
    session.pop('user_id', None)
    flash("Logged out")
    return redirect(url_for('home'))

@app.route('/election/<int:election_id>', methods=['GET', 'POST'])
@login_required
def election_detail(election_id):
    election = Election.query.get_or_404(election_id)
    user = User.query.get(session['user_id'])
    existing_vote = Vote.query.filter_by(user_id=user.id, election_id=election.id).first()

    if request.method == 'POST':
        # Ranked choices submitted as candidate ids in order
        ranked = request.form.getlist('ranked')
        # Validate ranked choices
        candidate_ids = {str(c.id) for c in election.candidates}
        if not ranked or any(c not in candidate_ids for c in ranked):
            flash("Invalid vote submission")
            return redirect(url_for('election_detail', election_id=election.id))
        # Save or update vote
        ranked_str = ','.join(ranked)
        if existing_vote:
            existing_vote.ranked_choices = ranked_str
        else:
            vote = Vote(user_id=user.id, election_id=election.id, ranked_choices=ranked_str)
            db.session.add(vote)
        db.session.commit()
        flash("Vote recorded")
        return redirect(url_for('election_detail', election_id=election.id))

    winner = runoff_winner(election)
    winner_name = winner.name if winner else "No winner yet"

    return render_template_string('''
    <!doctype html>
    <title>{{ election.name }}</title>
    <h1>{{ election.name }}</h1>
    <p>Winner: <strong>{{ winner_name }}</strong></p>
    {% if existing_vote %}
        <p>Your current vote ranking:</p>
        <ol>
        {% for cid in existing_vote.ranked_choices.split(',') %}
            <li>{{ candidates_dict[cid] }}</li>
        {% endfor %}
        </ol>
        <p>You can change your vote below.</p>
    {% else %}
        <p>You have not voted yet. Please rank the candidates below.</p>
    {% endif %}
    <form method="post">
      <p>Rank candidates by selecting their order:</p>
      <ul id="candidate-list" style="list-style:none; padding:0;">
      {% for candidate in election.candidates %}
        <li>
          <label>
            <input type="checkbox" name="candidate_checkbox" value="{{ candidate.id }}" onchange="updateRanking()">
            {{ candidate.name }}
          </label>
        </li>
      {% endfor %}
    ''')
