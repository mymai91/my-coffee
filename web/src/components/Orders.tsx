import { useOrders } from "../hooks";

const STATUS_EMOJI: Record<string, string> = {
  QUEUED: "🕐",
  GRINDING: "⚙️",
  BREWING: "☕",
  FROTHING: "🥛",
  READY: "✅",
};

export default function Orders() {
  const { data: orders = [], isLoading, error, refetch, isFetching } = useOrders();

  if (isLoading) return <div className="loading">Loading orders…</div>;
  if (error) return <div className="error">⚠️ {error.message}</div>;

  return (
    <div className="card">
      <div className="card-header">
        <h2>📦 Orders</h2>
        <button className="btn btn-small" onClick={() => refetch()} disabled={isFetching}>
          {isFetching ? "⏳ Loading…" : "🔄 Refresh"}
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
