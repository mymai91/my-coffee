import { useState } from "react";
import Menu from "./components/Menu";
import Orders from "./components/Orders";
import OrderForm from "./components/OrderForm";
import "./App.css";

type Tab = "menu" | "orders" | "order";

function App() {
  const [tab, setTab] = useState<Tab>("menu");

  return (
    <div className="app">
      <header className="header">
        <h1>☕ My Coffee Shop</h1>
        <p className="subtitle">Order your favorite brew</p>
      </header>

      <nav className="tabs">
        <button
          className={`tab ${tab === "menu" ? "active" : ""}`}
          onClick={() => setTab("menu")}
        >
          📋 Menu
        </button>
        <button
          className={`tab ${tab === "orders" ? "active" : ""}`}
          onClick={() => setTab("orders")}
        >
          📦 Orders
        </button>
        <button
          className={`tab ${tab === "order" ? "active" : ""}`}
          onClick={() => setTab("order")}
        >
          🛒 New Order
        </button>
      </nav>

      <main className="content">
        {tab === "menu" && <Menu />}
        {tab === "orders" && <Orders />}
        {tab === "order" && (
          <OrderForm onOrderCreated={() => setTab("orders")} />
        )}
      </main>
    </div>
  );
}

export default App;
