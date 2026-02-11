import { useState } from "react";
import { useCreateOrder } from "../hooks";

interface Props {
  onOrderCreated: () => void;
}

export default function OrderForm({ onOrderCreated }: Props) {
  const [name, setName] = useState("");
  const mutation = useCreateOrder();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    mutation.mutate(name.trim(), {
      onSuccess: () => {
        setName("");
        onOrderCreated();
      },
    });
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
          disabled={mutation.isPending}
          className="input"
        />
        <button
          type="submit"
          disabled={mutation.isPending || !name.trim()}
          className="btn"
        >
          {mutation.isPending ? "Ordering…" : "Order ☕"}
        </button>
      </form>
      {mutation.isSuccess && (
        <p className="message success">
          ✅ Order placed! ID: {mutation.data.orderId}
        </p>
      )}
      {mutation.isError && (
        <p className="message error">❌ {mutation.error.message}</p>
      )}
    </div>
  );
}
