from flask import Flask, request, redirect, url_for, render_template_string
from flask_sqlalchemy import SQLAlchemy

app = Flask(__name__)
app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite:///:memory:"
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False

db = SQLAlchemy(app)

# --------------------
# Database Models
# --------------------
class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)

class Vote(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey("user.id"), nullable=False)
    role = db.Column(db.String(50), nullable=False)
    nominee = db.Column(db.String(100), nullable=False)

# --------------------
# Initial Data
# --------------------
ROLES = {
    "President": ["Alice", "Bob"],
    "Vice-President": ["Charlie", "Diana"],
    "Scribe": ["Eve", "Frank"],
    "Registrar": ["Grace", "Henry"]
}

# --------------------
# Routes
# --------------------
@app.route("/", methods=["GET", "POST"])
def register():
    if request.method == "POST":
        name = request.form["name"]
        user = User(name=name)
        db.session.add(user)
        db.session.commit()
        return redirect(url_for("vote", user_id=user.id))

    return render_template_string("""
    <h2>User Registration</h2>
    <form method="post">
        <label>Name:</label><br>
        <input type="text" name="name" required><br><br>
        <button type="submit">Register</button>
    </form>
    """)

@app.route("/vote/<int:user_id>", methods=["GET", "POST"])
def vote(user_id):
    user = User.query.get_or_404(user_id)

    if request.method == "POST":
        for role in ROLES:
            nominee = request.form.get(role)
            if nominee:
                vote = Vote(user_id=user.id, role=role, nominee=nominee)
                db.session.add(vote)
        db.session.commit()
        return redirect(url_for("results"))

    return render_template_string("""
    <h2>Welcome, {{ user.name }}</h2>
    <form method="post">
        {% for role, nominees in roles.items() %}
            <h3>{{ role }}</h3>
            {% for nominee in nominees %}
                <input type="radio" name="{{ role }}" value="{{ nominee }}" required>
                {{ nominee }}<br>
            {% endfor %}
        {% endfor %}
        <br>
        <button type="submit">Submit Votes</button>
    </form>
    """, user=user, roles=ROLES)

@app.route("/results")
def results():
    results = {}
    for role in ROLES:
        votes = Vote.query.filter_by(role=role).all()
        tally = {}
        for v in votes:
            tally[v.nominee] = tally.get(v.nominee, 0) + 1
        results[role] = tally

    return render_template_string("""
    <h2>Election Results</h2>
    {% for role, tally in results.items() %}
        <h3>{{ role }}</h3>
        <ul>
        {% for nominee, count in tally.items() %}
            <li>{{ nominee }}: {{ count }} vote(s)</li>
        {% else %}
            <li>No votes yet</li>
        {% endfor %}
        </ul>
    {% endfor %}
    <a href="/">Register another user</a>
    """, results=results)

# --------------------
# App Startup
# --------------------
if __name__ == "__main__":
    with app.app_context():
        db.create_all()
    app.run(debug=True)
