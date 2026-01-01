from flask import Flask, request, redirect, url_for, render_template_string, jsonify
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.orm import Mapped

# --------------------
# App & Database
# --------------------
app = Flask(__name__)
app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite:///:memory:"
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False

db = SQLAlchemy(app)


# --------------------
# Model
# --------------------
class Todo(db.Model):
    id: Mapped[int] = db.Column(db.Integer, primary_key=True)
    title: Mapped[str] = db.Column(db.String(200), nullable=False)
    completed: Mapped[bool] = db.Column(db.Boolean, default=False)


with app.app_context():
    db.create_all()


# --------------------
# Embedded HTML (PWA)
# --------------------
HTML = """
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>TODO PWA</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="theme-color" content="#2c3e50">

    <link rel="manifest" href="/manifest.json">

    <style>
        body { font-family: Arial; margin: 30px; background: #f4f4f4; }
        h1 { color: #2c3e50; }
        form { margin-bottom: 20px; }
        input[type=text] { padding: 8px; width: 250px; }
        button { padding: 8px 12px; }
        ul { list-style: none; padding: 0; }
        li { background: #fff; margin-bottom: 8px; padding: 10px; display: flex; justify-content: space-between; }
        .done { text-decoration: line-through; color: gray; }
        a { margin-left: 8px; text-decoration: none; }
    </style>
</head>
<body>

<h1>📝 TODO PWA</h1>

<form method="post" action="/add">
    <input type="text" name="title" placeholder="New task..." required>
    <button>Add</button>
</form>

<ul>
{% for todo in todos %}
    <li>
        <span class="{{ 'done' if todo.completed else '' }}">
            {{ todo.title }}
        </span>
        <span>
            <a href="/toggle/{{ todo.id }}">✔</a>
            <a href="/delete/{{ todo.id }}">🗑</a>
        </span>
    </li>
{% endfor %}
</ul>

<script>
if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("/sw.js");
}
</script>

</body>
</html>
"""


# --------------------
# Routes (CRUD)
# --------------------
@app.route("/")
def index():
    todos = Todo.query.order_by(Todo.id.desc()).all()
    return render_template_string(HTML, todos=todos)


@app.route("/add", methods=["POST"])
def add():
    title = request.form["title"]
    todo = Todo(title=title)
    db.session.add(todo)
    db.session.commit()
    return redirect(url_for("index"))


@app.route("/toggle/<int:todo_id>")
def toggle(todo_id):
    todo = Todo.query.get_or_404(todo_id)
    todo.completed = not todo.completed
    db.session.commit()
    return redirect(url_for("index"))


@app.route("/delete/<int:todo_id>")
def delete(todo_id):
    todo = Todo.query.get_or_404(todo_id)
    db.session.delete(todo)
    db.session.commit()
    return redirect(url_for("index"))


# --------------------
# PWA Manifest
# --------------------
@app.route("/manifest.json")
def manifest():
    return jsonify({
        "name": "TODO PWA",
        "short_name": "TODO",
        "start_url": "/",
        "display": "standalone",
        "background_color": "#ffffff",
        "theme_color": "#2c3e50",
        "icons": []
    })


# --------------------
# Service Worker
# --------------------
@app.route("/sw.js")
def service_worker():
    js = """
self.addEventListener('install', event => {
    event.waitUntil(
        caches.open('todo-cache').then(cache => {
            return cache.addAll(['/']);
        })
    );
});

self.addEventListener('fetch', event => {
    event.respondWith(
        caches.match(event.request).then(response => {
            return response || fetch(event.request);
        })
    );
});
"""
    return app.response_class(js, mimetype="application/javascript")


# --------------------
# Run
# --------------------
if __name__ == "__main__":
    app.run(debug=True)
