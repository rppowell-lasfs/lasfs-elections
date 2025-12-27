from flask import Flask, request, redirect, url_for, render_template_string
from flask_sqlalchemy import SQLAlchemy

app = Flask(__name__)
app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite:///:memory:"
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False

db = SQLAlchemy(app)

# --------------------
# Models
# --------------------
class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)

class Election(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    title = db.Column(db.String(100), nullable=False)

class Role(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(50), nullable=False)
    election_id = db.Column(db.Integer, db.ForeignKey("election.id"))

class Nominee(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)
    role_id = db.Column(db.Integer, db.ForeignKey("role.id"))

class Vote(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer)
    nominee_id = db.Column(db.Integer)
    election_id = db.Column(db.Integer)

# --------------------
# User Routes
# --------------------
@app.route("/", methods=["GET", "POST"])
def register():
    if request.method == "POST":
        user = User(name=request.form["name"])
        db.session.add(user)
        db.session.commit()
        return redirect(url_for("select_election", user_id=user.id))

    return render_template_string("""
    <h2>User Registration</h2>
    <form method="post">
        <input name="name" placeholder="Your name" required>
        <button>Register</button>
    </form>
    <br>
    <a href="/admin">Admin Dashboard</a>
    """)

@app.route("/elections/<int:user_id>")
def select_election(user_id):
    elections = Election.query.all()
    return render_template_string("""
    <h2>Select Election</h2>
    {% for e in elections %}
        <a href="/vote/{{ user_id }}/{{ e.id }}">{{ e.title }}</a><br>
    {% endfor %}
    """, elections=elections, user_id=user_id)

@app.route("/vote/<int:user_id>/<int:election_id>", methods=["GET", "POST"])
def vote(user_id, election_id):
    roles = Role.query.filter_by(election_id=election_id).all()

    if request.method == "POST":
        for role in roles:
            nominee_id = request.form.get(str(role.id))
            if nominee_id:
                vote = Vote(
                    user_id=user_id,
                    nominee_id=int(nominee_id),
                    election_id=election_id
                )
                db.session.add(vote)
        db.session.commit()
        return redirect(url_for("results", election_id=election_id))

    return render_template_string("""
    <h2>Voting</h2>
    <form method="post">
        {% for role in roles %}
            <h3>{{ role.name }}</h3>
            {% for n in role.nominees %}
                <input type="radio" name="{{ role.id }}" value="{{ n.id }}" required>
                {{ n.name }}<br>
            {% endfor %}
        {% endfor %}
        <br>
        <button>Submit Vote</button>
    </form>
    """, roles=roles)

# --------------------
# Admin Routes
# --------------------
@app.route("/admin")
def admin():
    elections = Election.query.all()
    return render_template_string("""
    <h2>Admin Dashboard</h2>
    <a href="/admin/election/create">Create Election</a><br><br>
    {% for e in elections %}
        <b>{{ e.title }}</b>
        | <a href="/admin/election/{{ e.id }}">Manage</a>
        | <a href="/results/{{ e.id }}">Results</a>
        <br>
    {% endfor %}
    """)

@app.route("/admin/election/create", methods=["GET", "POST"])
def create_election():
    if request.method == "POST":
        election = Election(title=request.form["title"])
        db.session.add(election)
        db.session.commit()
        return redirect(url_for("admin"))

    return render_template_string("""
    <h2>Create Election</h2>
    <form method="post">
        <input name="title" placeholder="Election title" required>
        <button>Create</button>
    </form>
    """)

@app.route("/admin/election/<int:election_id>", methods=["GET", "POST"])
def manage_election(election_id):
    election = Election.query.get_or_404(election_id)

    if request.method == "POST":
        role = Role(
            name=request.form["role"],
            election_id=election.id
        )
        db.session.add(role)
        db.session.commit()

    roles = Role.query.filter_by(election_id=election.id).all()

    return render_template_string("""
    <h2>Manage {{ election.title }}</h2>

    <h3>Add Role</h3>
    <form method="post">
        <input name="role" placeholder="Role name" required>
        <button>Add</button>
    </form>

    <h3>Roles</h3>
    {% for r in roles %}
        {{ r.name }} |
        <a href="/admin/role/{{ r.id }}">Add Nominees</a><br>
    {% endfor %}
    """, election=election, roles=roles)

@app.route("/admin/role/<int:role_id>", methods=["GET", "POST"])
def manage_role(role_id):
    role = Role.query.get_or_404(role_id)

    if request.method == "POST":
        nominee = Nominee(
            name=request.form["name"],
            role_id=role.id
        )
        db.session.add(nominee)
        db.session.commit()

    nominees = Nominee.query.filter_by(role_id=role.id).all()

    return render_template_string("""
    <h2>Nominees for {{ role.name }}</h2>
    <form method="post">
        <input name="name" placeholder="Nominee name" required>
        <button>Add</button>
    </form>

    <ul>
        {% for n in nominees %}
            <li>{{ n.name }}</li>
        {% endfor %}
    </ul>
    """ , role=role, nominees=nominees)

# --------------------
# Results
# --------------------
@app.route("/results/<int:election_id>")
def results(election_id):
    election = Election.query.get_or_404(election_id)
    roles = Role.query.filter_by(election_id=election.id).all()

    data = []
    for role in roles:
        tally = {}
        for n in Nominee.query.filter_by(role_id=role.id):
            count = Vote.query.filter_by(
                nominee_id=n.id,
                election_id=election.id
            ).count()
            tally[n.name] = count
        data.append((role.name, tally))

    return render_template_string("""
    <h2>Results: {{ election.title }}</h2>
    {% for role, tally in data %}
        <h3>{{ role }}</h3>
        <ul>
            {% for name, count in tally.items() %}
                <li>{{ name }} — {{ count }}</li>
            {% endfor %}
        </ul>
    {% endfor %}
    """, election=election, data=data)

# --------------------
# Start App
# --------------------
if __name__ == "__main__":
    with app.app_context():
        db.create_all()
    app.run(debug=True)
