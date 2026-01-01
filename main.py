from enum import Enum
from flask import Flask, render_template_string, request, redirect, url_for, session
from sqlalchemy.orm import Mapped, mapped_column
from flask_sqlalchemy import SQLAlchemy
from flask_wtf import FlaskForm
from wtforms import StringField, SelectField, SubmitField, PasswordField
from wtforms.validators import DataRequired, Length
from werkzeug.security import generate_password_hash, check_password_hash
from typing import List
import logging

def create_app(testing=False):
    app = Flask(__name__)
    app.logger.setLevel(logging.DEBUG)
    app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///:memory:'
    app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
    app.config['SECRET_KEY'] = 'supersecretkey'
    app.config['TESTING'] = testing

    return app

app = create_app()

db = SQLAlchemy(app)

# --------------------
# Models
# --------------------
class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    password_hash = db.Column(db.String(128), nullable=False)

    def set_password(self, password):
        self.password_hash = generate_password_hash(password)

    def check_password(self, password):
        return check_password_hash(self.password_hash, password)


class ElectionStatus(Enum):
    OPEN = "Open"
    CLOSED = "Closed"

class Election(db.Model):
    id: Mapped[int] = mapped_column('id', primary_key=True)
    name: Mapped[str] = mapped_column('name', nullable=False)
    status: Mapped[str] = mapped_column('status', nullable=False, default="Closed")
    # status = db.Column(db.Enum(ElectionStatus), nullable=False)
    positions: Mapped[List['Position']] = db.relationship(back_populates='election')
    # positions = db.relationship('Position', back_populates='election')

class Position(db.Model):
    id: Mapped[int] = mapped_column('id', primary_key=True)
    name: Mapped[str] = mapped_column('name', nullable=False)
    election_id: Mapped[int] = mapped_column(db.ForeignKey("election.id"))
    election: Mapped[Election] = db.relationship(back_populates='positions')
    nominees: Mapped[List['Nominee']] = db.relationship(back_populates='position')
    # election = db.relationship('Election', back_populates='positions')
    # nominees: Mapped[list['Nominee']] = db.relationship('Nominee',
    #     back_populates='nominee'
    #     # backref='nominee'
    #     # lazy='dynamic'
    # )
    # # election = db.relationship('Election', backref='')

class Nominee(db.Model):
    id: Mapped[int] = mapped_column('id', primary_key=True)
    name: Mapped[str] = mapped_column('name', nullable=False)
    position_id: Mapped[int] = mapped_column('position_id', db.ForeignKey("position.id"))
    position: Mapped[Position] = db.relationship(back_populates='nominees')

class Vote(db.Model):
    id: Mapped[int] = db.Column(db.Integer, primary_key=True)
    user_id: Mapped[int] = db.Column(db.Integer)
    nominee_id: Mapped[int] = db.Column(db.Integer)


TEST_ELECTIONS = [
    {'name':'test1', 'positions':[
        {'name': "President", 'nominees': ["Matthew", "Michelle", "Karl"]}
    ]}
]

TEST_USERS = [
    {'username':'test1', 'password':'test1'}
]

def add_TEST_ELECTIONS(db):
    for TEST_ELECTION in TEST_ELECTIONS:
        election = Election(name=TEST_ELECTION['name'])
        db.session.add(election)
        db.session.commit()
        # app.logger.debug(["add_TEST_ELECTIONS TEST_ELECTION", election.id ]+TEST_ELECTION['positions'])
        for TEST_POSITION in TEST_ELECTION['positions']:
            position = Position(name=TEST_POSITION['name'], election_id=election.id)
            db.session.add(position)
            db.session.commit()
            # app.logger.debug(["add_TEST_ELECTIONS TEST_POSITION"] + [position.id])
            # app.logger.debug(["add_TEST_ELECTIONS TEST_POSITION"] + TEST_POSITION['nominees'])
            for TEST_NOMINEE in TEST_POSITION['nominees']:
                nominee = Nominee(name=TEST_NOMINEE, position_id = position.id)
                db.session.add(nominee)
                db.session.commit()
    db.session.commit()

def add_TEST_USERS(db):
    for TEST_USER in TEST_USERS:
        new_user = User(username=TEST_USER['username'])
        new_user.set_password(TEST_USER['password'])
        db.session.add(new_user)
        db.session.commit()
    db.session.commit()


# home_html = """
# <!doctype html>
# <title>Home</title>
# <h1>Welcome {% if user %}{{ user }}{% else %}Guest{% endif %}</h1>
# {% if user %}
#   <a href="{{ url_for('logout') }}">Logout</a>
# {% else %}
#   <a href="{{ url_for('login') }}">Login</a> | <a href="{{ url_for('register') }}">Register</a>
# {% endif %}
# """

# register_html = """
# <!doctype html>
# <title>Register</title>
# <h1>Register</h1>
# <form method="post">
#   Username: <input type="text" name="username" required><br>
#   Password: <input type="password" name="password" required><br>
#   <input type="submit" value="Register">
# </form>
# {% if error %}
# <p style="color:red;">{{ error }}</p>
# {% endif %}
# <a href="{{ url_for('home') }}">Home</a>
# """

# register_html = """
# <!doctype html>
# <title>Register</title>
# <h1>Register</h1>
# <form method="post" action="/login">
#   Username: <input type="text" name="username" required><br>
#   Password: <input type="password" name="password" required><br>
#   <input type="submit" value="Register">
# </form>
# <hr>
# <form method="POST" action="/login">
#     {{ form.hidden_tag() }}
#     <p>
#         {{ form.username.label }}
#         {{ form.username }}
#         {% for error in form.username.errors %}
#             <span style="color: red;">[{{ error }}]</span>
#         {% endfor %}
#     </p>
#     <p>
#         {{ form.password.label }}<br>
#         {{ form.password }}
#         {% for error in form.password.errors %}
#             <span style="color: red;">[{{ error }}]</span>
#         {% endfor %}
#     </p>
#     <p>{{ form.submit() }}</p>
# </form>
# {% if error %}
# <p style="color:red;">{{ error }}</p>
# {% endif %}
# <a href="{{ url_for('home') }}">Home</a>
# """

# login_html = """
# <!doctype html>
# <title>Login</title>
# <h1>Login</h1>
# <form method="post">
#   Username: <input type="text" name="username" required><br>
#   Password: <input type="password" name="password" required><br>
#   <input type="submit" value="Login">
# </form>
# {% if error %}
# <p style="color:red;">{{ error }}</p>
# {% endif %}
# <a href="{{ url_for('home') }}">Home</a>
# """

# login_html = """
# <!doctype html>
# <title>Login</title>
# <h1>Login</h1>
# <form method="POST" action="/login">
#     {{ form.hidden_tag() }}
#     <p>
#         {{ form.username.label }}<br>
#         {{ form.username }}
#         {% for error in form.username.errors %}
#             <span style="color: red;">[{{ error }}]</span>
#         {% endfor %}
#     </p>
#     <p>
#         {{ form.password.label }}<br>
#         {{ form.password }}
#         {% for error in form.password.errors %}
#             <span style="color: red;">[{{ error }}]</span>
#         {% endfor %}
#     </p>
#     <p>{{ form.submit() }}</p>
# </form>
# {% if error %}
# <p style="color:red;">{{ error }}</p>
# {% endif %}
# <a href="{{ url_for('home') }}">Home</a>
# """


class ElectionForm(FlaskForm):
    name = StringField('Name', validators=[DataRequired()], default="Election name")

class NewUserForm(FlaskForm):
    username = StringField('Username', validators=[DataRequired(), Length(min=4, max=25)])
    password = PasswordField('Password', validators=[DataRequired(), Length(min=4)])
    submit = SubmitField('Create')

class LoginForm(FlaskForm):
    username = StringField('Username', validators=[DataRequired()])
    password = PasswordField('Password', validators=[DataRequired()])
    submit = SubmitField('Sign In')

@app.route('/')
def home():
    user = session.get('user')
    elections = Election.query.all() if user else []
    app.logger.debug(elections)
    return render_template_string("""
<!doctype html>
<title>Home</title>
<h1>Welcome {% if user %}{{ user }}{% else %}Guest{% endif %}</h1>
{% if user %}
  <h2>Election Dashboard</h2>
    {% for e in elections %}
      <b><a href="/election/{{ e.id }}">{{ e.name }}</a></b>
      | <a href="/results/{{ e.id }}">Results</a>
      <br>
  {% endfor %}
  </hr>
  <a href="{{ url_for('logout') }}">Logout</a>

{% else %}
  <a href="{{ url_for('login') }}">Login</a> | <a href="{{ url_for('register') }}">Register</a>
{% endif %}
""", user=user, elections=elections)

@app.route('/register', methods=['GET', 'POST'])
def register():
    error = None
    form = NewUserForm()
    if form.validate_on_submit():
        username = form.username.data.strip()
        password = form.password.data
        if User.query.filter_by(username=username).first():
            error = 'Username already exists.'
        else:
            new_user = User(username=username)
            new_user.set_password(password)
            db.session.add(new_user)
            db.session.commit()
            session['user'] = username
            return redirect(url_for('home'))
    return render_template_string("""
<!doctype html>
<title>Register</title>
<h1>Register</h1>
<form method="POST" action="/register">
    {{ form.hidden_tag() }}
    <p>
        {{ form.username.label }}:
        {{ form.username }}
        {% for error in form.username.errors %}
            <span style="color: red;">[{{ error }}]</span>
        {% endfor %}
    </p>
    <p>
        {{ form.password.label }}:
        {{ form.password }}
        {% for error in form.password.errors %}
            <span style="color: red;">[{{ error }}]</span>
        {% endfor %}
    </p>
    <p>{{ form.submit() }}</p>
</form>
{% if error %}
<p style="color:red;">{{ error }}</p>
{% endif %}
<a href="{{ url_for('home') }}">Home</a>
""", form=form, error=error)

@app.route('/login', methods=['GET', 'POST'])
def login():
    error = None
    form = LoginForm()
    if form.validate_on_submit():
        username = form.username.data.strip()
        password = form.password.data
        user = User.query.filter_by(username=username).first()
        if user and user.check_password(password):
            session['user'] = username
            return redirect(url_for('home'))
        else:
            error = 'Invalid username or password.'
    return render_template_string("""
<!doctype html>
<title>Login</title>
<h1>Login</h1>
<form method="POST" action="/login">
    {{ form.hidden_tag() }}
    <p>
        {{ form.username.label }}<br>
        {{ form.username }}
        {% for error in form.username.errors %}
            <span style="color: red;">[{{ error }}]</span>
        {% endfor %}
    </p>
    <p>
        {{ form.password.label }}<br>
        {{ form.password }}
        {% for error in form.password.errors %}
            <span style="color: red;">[{{ error }}]</span>
        {% endfor %}
    </p>
    <p>{{ form.submit() }}</p>
</form>
{% if error %}
<p style="color:red;">{{ error }}</p>
{% endif %}
<a href="{{ url_for('home') }}">Home</a>
""", form=form, error=error)

@app.route('/logout')
def logout():
    session.pop('user', None)
    return redirect(url_for('home'))


# --------------------
# Admin Routes
# --------------------
@app.route("/admin")
def admin():
    elections = Election.query.all()
    app.logger.debug(elections)
    return render_template_string("""
<h2>Admin Dashboard</h2>
<a href="/admin/election/create">Create Election</a><br><br>
{% for e in elections %}
    <b><a href="/admin/election/{{ e.id }}">{{ e.name }} ({{ e.status}})</a></b>
    | <a href="/results/{{ e.id }}">Results</a>
    <br>
{% endfor %}
""", elections=elections)

@app.route("/admin/election/create", methods=["GET", "POST"])
def create_election():
    if request.method == "POST":
        election = Election(name=request.form["name"])
        db.session.add(election)
        db.session.commit()
        return redirect(url_for("admin"))
    else:
        election = Election()
    return render_template_string("""
<h2>Create Election</h2>
<form method="post">
    <input name="name" placeholder="Election name" required>
    <select name="status" id="status">
        {% for status in statuses %}
            <option value="{{ status.value }}"
                {% if status == election.status %}selected{% endif %}>
                {{ status.value }}
            </option>
        {% endfor %}
    </select>        
    <button>Create</button>
</form>
""", election=election, statuses=ElectionStatus)

@app.route("/admin/election/<int:election_id>", methods=["GET", "POST"])
def manage_election(election_id):
    election = Election.query.get_or_404(election_id)

    if request.method == "POST":
        position = Position(
            name=request.form["position"],
            election_id=election.id
        )
        db.session.add(position)
        db.session.commit()

    positions = Position.query.filter_by(election_id=election.id).all()

    app.logger.debug([[position.name, position.nominees] for position in positions])
    return render_template_string("""
<h2>Manage Election {{ election.name }} ({{ election.status }})</h2>

<h3>Positions</h3>
{% for p in positions %}
    {{ p.name }} |
    {% for n in p.nominees %}
        {{ n.name }} |
    {% endfor %}
    <a href="/admin/position/{{ p.id }}">Add Nominees</a><br>
{% endfor %}

<h3>Add Position</h3>
<form method="post">
    <input name="position" placeholder="Position name" required>
    <button>Add</button>
</form>
<a href="/admin">Admin Dashboard</a><br>
""", election=election, positions=positions)

@app.route("/admin/position/<int:position_id>", methods=["GET", "POST"])
def manage_role(position_id):
    position = Position.query.get_or_404(position_id)

    if request.method == "POST":
        nominee = Nominee(
            name=request.form["name"],
            position_id=position.id
        )
        db.session.add(nominee)
        db.session.commit()

    nominees = Nominee.query.filter_by(position_id=position.id).all()

    return render_template_string("""
<h2>Nominees for {{ position.name }}</h2>
<form method="post">
    <input name="name" placeholder="Nominee name" required>
    <button>Add</button>
</form>
<ul>
    {% for n in nominees %}
        <li>{{ n.name }}</li>
    {% endfor %}
</ul>
<a href="/admin/election/{{ position.election_id }}">Back to {{ position.election.name }} </a><br>
""" , position=position, nominees=nominees)


if __name__ == '__main__':

    with app.app_context():
        db.create_all()
        db.session.commit()
        add_TEST_ELECTIONS(db)
        add_TEST_USERS(db)

    app.run(debug=True)