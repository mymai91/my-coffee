import { useState } from "react";
import { createOrder } from "../api";

interface Props {
  onOrderCreated: () => void;
}

export default function OrderForm({ onOrderCreated }: Props) {
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [isError, setIsError] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setLoading(true);
    setMessage("");
    setIsError(false);

    try {
      const res = await createOrder(name.trim());
      setMessage(`✅ Order placed! ID: ${res.orderId}`);
      setName("");
      onOrderCreated();
    } catch (err: unknown) {
      setIsError(true);
      setMessage(
        `❌ ${err instanceof Error ? err.message : "Failed to place order"}`
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card">
      <h2>🛒 Place an Order</h2>
      <form onSubmit={handleSubmit} className="order-form">
        <input
          type="text"
          placeholder="Enter drink name (e.g. Latte)"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={loading}
          className="input"
        />
        <button type="submit" disabled={loading || !name.trim()} className="btn">
          {loading ? "Ordering…" : "Order ☕"}
        </button>
      </form>
      {message && (
        <p className={`message ${isError ? "error" : "success"}`}>{message}</p>
      )}
    </div>
  );
}
