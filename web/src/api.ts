const API_BASE = "http://localhost:9000";

export interface MenuItem {
  name: string;
  description: string;
  price: number;
}

export interface Order {
  orderId: string;
  menuItemName: string;
  status: string;
}

export async function fetchMenu(): Promise<MenuItem[]> {
  const res = await fetch(`${API_BASE}/api/menu`);
  if (!res.ok) throw new Error("Failed to fetch menu");
  return res.json();
}

export async function fetchOrders(): Promise<Order[]> {
  const res = await fetch(`${API_BASE}/api/orders`);
  if (!res.ok) throw new Error("Failed to fetch orders");
  return res.json();
}

export async function createOrder(
  menuItemName: string
): Promise<{ orderId: string }> {
  const res = await fetch(`${API_BASE}/api/orders`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ menuItemName }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || "Failed to create order");
  }
  return res.json();
}
