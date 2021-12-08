import "water.css";
import "./App.css";

import { useEffect, useState } from "react";

function App() {
  const [ranges, setRanges] = useState([]);
  const [username, setUsername] = useState("");
  const [inputUsername, setInputUsername] = useState("");
  const [token, setToken] = useState("");
  const [category, setCategory] = useState("");
  const [limit, setLimit] = useState(1);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState("");

  function startAPI() {
    fetch("/api/start", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        token: token,
      }),
    }).then((response) => {
      if (!response.ok) {
        response.text().then((text) => {
          alert(text);
        });
      }
    });
  }

  function stopAPI() {
    fetch("/api/stop", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        token: token,
      }),
    }).then((response) => {
      if (!response.ok) {
        response.text().then((text) => {
          alert(text);
        });
      }
    });
  }

  function setLimitAPI(limit) {
    fetch(`/api/setlimit`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        token: token,
        limit: parseInt(limit),
      }),
    }).then((response) => {
      if (!response.ok) {
        response.text().then((text) => {
          alert(text);
        });
      }
    });
  }

  function refreshRange() {
    fetch("/api/range", {
    }).then((res) => {
      if (!res.ok) {
        res.text().then((text) => {
          alert(text);
        });
      }
      res.json().then((json) => {
        setRanges(json);
      });
    });
  }

  function addRange(category, start, end, token) {
    console.log(convertNumberToDate(start));
    fetch("/api/range/add", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        token: token,
        range: {
          category: category,
          start: start,
          end: end,
        },
      }),
    }).then(refreshRange);
  }

  function setRangeUser(id, user) {
    fetch("/api/range/setuser", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        range: {
          id: id,
          user: user,
        },
        token: token,
      }),
    })
      .then((res) => {
        if (!res.ok) {
          // alert respond text
          res.text().then((text) => alert(text));
        }
      })
      .then(refreshRange);
  }

  function clearRange(id) {
    fetch("/api/range/clear", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: id,
        token: token,
      }),
    })
      .then((res) => {
        if (!res.ok) {
          // alert respond text
          res.text().then((text) => alert(text));
        }
      })
      .then(refreshRange);
  }

  function deleteRange(id) {
    fetch("/api/range/delete", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: id,
        token: token,
      }),
    })
      .then((res) => {
        if (!res.ok) {
          // alert respond text
          res.text().then((text) => alert(text));
        }
      })
      .then(refreshRange);
  }

  function getCategories() {
    fetch("/api/category")
      .then((res) => res.json())
      .then((data) => {
        setCategories(data);
        setSelectedCategory(data[0].name);
      });
  }

  useEffect(() => {
    refreshRange();
  }, []);

  return (
    <div className="App">
      <h1>Hi, {username ? username : "你的大名？"}</h1>

      {username === "" && (
        <span>
          <input
            type="text"
            value={inputUsername}
            placeholder="输入姓名"
            onChange={(e) => setInputUsername(e.target.value)}
          />
          <button onClick={() => setUsername(inputUsername)}>Set</button>
        </span>
      )}

      {username === "admin" && (
        <div>
          <input
            type="text"
            value={token}
            placeholder="token"
            onChange={(e) => setToken(e.target.value)}
          />
          <input
            type="text"
            value={category}
            placeholder="category"
            onChange={(e) => setCategory(e.target.value)}
          />
          <input id="start" type="datetime-local" />
          <input id="end" type="datetime-local" />
          <button
            onClick={() =>
              addRange(
                category,
                document.getElementById("start").valueAsNumber,
                document.getElementById("end").valueAsNumber,
                token
              )
            }
          >
            Add
          </button>
          <input
            type="number"
            value={limit}
            placeholder="limit"
            onChange={(e) => setLimit(e.target.value)}
          />
          <button onClick={() => setLimitAPI(limit)}>Set Limit</button>
          <button onClick={() => startAPI()}>Start</button>
          <button onClick={() => stopAPI()}>Stop</button>
        </div>
      )}

      {username && (
        <div>
          <h2>时间段</h2>
          <p>找到想要的时间段，点击按钮即可。手快有手慢冇，Good luck!</p>

          <button onClick={() => refreshRange()}>刷新</button>
          <ul>
            {ranges.map((r) => (
              <li key={r.id}>
                <button
                  onClick={() => setRangeUser(r.id, username)}
                  disabled={r.user}
                >
                  {r.category} {convertNumberToDate(r.start, r.end)} {r.user}
                </button>
                {username === "admin" && (
                  <div>
                    <button onClick={() => clearRange(r.id)}>Clear</button>
                    <button onClick={() => deleteRange(r.id)}>Delete</button>
                    <input
                      type="text"
                      value={r.user}
                      onChange={(e) => setRangeUser(r.id, e.target.value)}
                    />
                  </div>
                )}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

function convertNumberToDate(startNumber, endNumber) {
  const start = new Date(startNumber);
  const end = new Date(endNumber);
  return `${
    start.getMonth() + 1
  }月${start.getDate()}日 ${start.getHours()}-${end.getHours()}`;
}

export default App;
