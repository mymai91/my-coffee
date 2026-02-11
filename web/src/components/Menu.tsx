import { useEffect, useState } from "react";
import { fetchMenu } from "../api";
import type { MenuItem } from "../api";

export default function Menu() {
  const [items, setItems] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchMenu()
      .then(setItems)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="loading">Loading menu…</div>;
  if (error) return <div className="error">⚠️ {error}</div>;

  return (
    <div className="card">
      <h2>📋 Menu</h2>
      {items.length === 0 ? (
        <p className="empty">No menu items available.</p>
      ) : (
        <div className="menu-grid">
          {items.map((item) => (
            <div key={item.name} className="menu-item">
              <div className="menu-item-header">
                <span className="menu-item-name">{item.name}</span>
                <span className="menu-item-price">
                  ${item.price.toFixed(2)}
                </span>
              </div>
              {item.description && (
                <p className="menu-item-desc">{item.description}</p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
