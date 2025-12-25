from flask import Flask, render_template_string, request, redirect, url_for, session
from flask_sqlalchemy import SQLAlchemy
from werkzeug.security import generate_password_hash, check_password_hash

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///:memory:'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
app.config['SECRET_KEY'] = 'supersecretkey'
db = SQLAlchemy(app)

class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    password_hash = db.Column(db.String(128), nullable=False)

    def set_password(self, password):
        self.password_hash = generate_password_hash(password)

    def check_password(self, password):
        return check_password_hash(self.password_hash, password)


ELECTIONS = [
    {'id':'1', 'position': "President", 'nominees': ["Matthew", "Michelle", "Karl"]}
]

home_html = """
<!doctype html>
<title>Home</title>
<h1>Welcome {% if user %}{{ user }}{% else %}Guest{% endif %}</h1>
{% if user %}
  <a href="{{ url_for('logout') }}">Logout</a>
{% else %}
  <a href="{{ url_for('login') }}">Login</a> | <a href="{{ url_for('register') }}">Register</a>
{% endif %}
"""

register_html = """
<!doctype html>
<title>Register</title>
<h1>Register</h1>
<form method="post">
  Username: <input type="text" name="username" required><br>
  Password: <input type="password" name="password" required><br>
  <input type="submit" value="Register">
</form>
{% if error %}
<p style="color:red;">{{ error }}</p>
{% endif %}
<a href="{{ url_for('home') }}">Home</a>
"""

login_html = """
<!doctype html>
<title>Login</title>
<h1>Login</h1>
<form method="post">
  Username: <input type="text" name="username" required><br>
  Password: <input type="password" name="password" required><br>
  <input type="submit" value="Login">
</form>
{% if error %}
<p style="color:red;">{{ error }}</p>
{% endif %}
<a href="{{ url_for('home') }}">Home</a>
"""

@app.route('/')
def home():
    user = session.get('user')
    return render_template_string(home_html, user=user)

@app.route('/register', methods=['GET', 'POST'])
def register():
    error = None
    if request.method == 'POST':
        username = request.form['username'].strip()
        password = request.form['password']
        if User.query.filter_by(username=username).first():
            error = 'Username already exists.'
        else:
            new_user = User(username=username)
            new_user.set_password(password)
            db.session.add(new_user)
            db.session.commit()
            session['user'] = username
            return redirect(url_for('home'))
    return render_template_string(register_html, error=error)

@app.route('/login', methods=['GET', 'POST'])
def login():
    error = None
    if request.method == 'POST':
        username = request.form['username'].strip()
        password = request.form['password']
        user = User.query.filter_by(username=username).first()
        if user and user.check_password(password):
            session['user'] = username
            return redirect(url_for('home'))
        else:
            error = 'Invalid username or password.'
    return render_template_string(login_html, error=error)

@app.route('/logout')
def logout():
    session.pop('user', None)
    return redirect(url_for('home'))

if __name__ == '__main__':
    with app.app_context():
        db.create_all()
    app.run(debug=True)