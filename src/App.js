import React, { useState, useEffect } from "react";
import API, { setAuthToken } from "./api";
import "./App.css"; // <-- IMPORT CSS

function App() {
  const [token, setToken] = useState("");
  const [view, setView] = useState("login");
  const [items, setItems] = useState([]);
  const [currentCartId, setCurrentCartId] = useState(null);

  const handleLogin = async (e) => {
    e.preventDefault();
    const username = e.target.username.value;
    const password = e.target.password.value;
    
    try {
      const res = await API.post("/users/login", { username, password });
      const t = res.data.token;
      setToken(t);
      setAuthToken(t);
      setView("items");
    } catch (err) {
      alert("Invalid username/password");
    }
  };

  useEffect(() => {
    const fetchItems = async () => {
      try {
        const res = await API.get("/items");
        setItems(res.data);
      } catch (err) {
        console.error(err);
      }
    };
    if (view === "items") fetchItems();
  }, [view]);

  const handleAddToCart = async (itemId) => {
    try {
      const res = await API.post("/carts", { item_id: itemId });
      setCurrentCartId(res.data.cart_id);
      alert(`Added item: ${itemId}`);
    } catch (err) {
      alert("Login required to continue");
    }
  };

  const handleShowCart = async () => {
    const res = await API.get("/carts");
    let msg = "Cart items:\n";
    res.data.forEach(cart => {
      (cart.items || []).forEach(i => {
        msg += `Cart: ${i.cart_id} Item: ${i.item_id}\n`;
      });
    });
    alert(msg);
  };

  const handleShowOrders = async () => {
    const res = await API.get("/orders");
    if (res.data.length === 0) return alert("No orders yet");
    alert(res.data.map(o => `Order ${o.id}`).join("\n"));
  };

  const handleCheckout = async () => {
    if (!currentCartId) return alert("Add items before checkout");
    await API.post("/orders", { cart_id: currentCartId });
    alert("Order placed!");
    setCurrentCartId(null);
  };

  /* -------- UI -------- */

  if (view === "login") {
    return (
      <div className="app-container">
        <h1 className="app-title">Shopping Login</h1>
        <form onSubmit={handleLogin}>
          <label>Username</label>
          <input name="username" />

          <label>Password</label>
          <input type="password" name="password" />

          <button type="submit">Login</button>
        </form>
      </div>
    );
  }

  return (
    <div className="app-container">
      <h1 className="app-title">Shopping Items</h1>

      <div className="top-buttons">
        <button onClick={handleCheckout}>Checkout</button>
        <button onClick={handleShowCart}>Cart</button>
        <button onClick={handleShowOrders}>Orders</button>
      </div>

      <ul className="items-list">
        {items.map(item => (
          <li key={item.id} onClick={() => handleAddToCart(item.id)}>
            {item.name} ({item.status})
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
