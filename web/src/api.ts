const BREW_BASE = "http://localhost:50051";
const MENU_BASE = "http://localhost:50052";

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

// Helper to call Connect RPC endpoints with JSON
async function connectFetch<T>(baseUrl: string, method: string, body: object = {}): Promise<T> {
  const res = await fetch(`${baseUrl}/${method}`, {
    method: "POST",  // Connect RPC always uses POST
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => null);
    throw new Error(err?.message || `Failed to call ${method}`);
  }
  return res.json();
}

export async function fetchMenu(): Promise<MenuItem[]> {
  const resp = await connectFetch<{ items: MenuItem[] }>(
    MENU_BASE, "menu.MenuService/GetMenu"
  );
  return resp.items ?? [];
}

export async function fetchOrders(): Promise<Order[]> {
  const resp = await connectFetch<{ orders: Order[] }>(
    BREW_BASE, "brew.BrewService/ListOrders"
  );
  return resp.orders ?? [];
}

export async function createOrder(
  menuItemName: string
): Promise<{ orderId: string }> {
  return connectFetch<{ orderId: string }>(
    BREW_BASE, "brew.BrewService/OrderDrink", { menuItemName }
  );
}

export async function getOrder(orderId: string): Promise<Order> {
  const resp = await connectFetch<{ order: Order }>(
    BREW_BASE, "brew.BrewService/GetOrder", { orderId }
  );
  return resp.order;
}

export async function updateOrderStatus(
  orderId: string, status: string
): Promise<Order> {
  // Map status string to enum value
  const statusMap: Record<string, number> = {
    QUEUED: 1,
    GRINDING: 2,
    BREWING: 3,
    FROTHING: 4,
    READY: 5,
  };
  const resp = await connectFetch<{ order: Order }>(
    BREW_BASE, "brew.BrewService/UpdateOrderStatus",
    { orderId, status: statusMap[status] ?? 0 }
  );
  return resp.order;
}

export async function deleteOrder(orderId: string): Promise<{ success: boolean }> {
  return connectFetch<{ success: boolean }>(
    BREW_BASE, "brew.BrewService/DeleteOrder", { orderId }
  );
}
