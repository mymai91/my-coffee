import { useEffect, useState } from "react";
import { fetchOrders } from "../api";
import type { Order } from "../api";

const STATUS_EMOJI: Record<string, string> = {
  QUEUED: "🕐",
  GRINDING: "⚙️",
  BREWING: "☕",
  FROTHING: "🥛",
  READY: "✅",
};

export default function Orders() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const load = () => {
    setLoading(true);
    fetchOrders()
      .then(setOrders)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, []);

  if (loading) return <div className="loading">Loading orders…</div>;
  if (error) return <div className="error">⚠️ {error}</div>;

  return (
    <div className="card">
      <div className="card-header">
        <h2>📦 Orders</h2>
        <button className="btn btn-small" onClick={load}>
          🔄 Refresh
        </button>
      </div>
      {orders.length === 0 ? (
        <p className="empty">No orders yet. Place one!</p>
      ) : (
        <table className="orders-table">
          <thead>
            <tr>
              <th>#</th>
              <th>Drink</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {orders.map((order, i) => (
              <tr key={order.orderId}>
                <td className="order-num">{i + 1}</td>
                <td>{order.menuItemName}</td>
                <td>
                  <span className={`status status-${order.status.toLowerCase()}`}>
                    {STATUS_EMOJI[order.status] ?? "❓"} {order.status}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
