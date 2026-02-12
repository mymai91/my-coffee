import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchMenu,
  fetchOrders,
  createOrder,
  getOrder,
  updateOrderStatus,
  deleteOrder,
} from "./api";

// Query key constants — avoids typos and makes invalidation easy
export const queryKeys = {
  menu: ["menu"] as const,
  orders: ["orders"] as const,
  order: (id: string) => ["order", id] as const,
};

/**
 * Fetches the coffee menu.
 * Cached for 5 minutes — menu rarely changes.
 */
export function useMenu() {
  return useQuery({
    queryKey: queryKeys.menu,
    queryFn: fetchMenu,
    staleTime: 5 * 60 * 1000,
  });
}

/**
 * Fetches all orders.
 * Refetches every 10s so the user sees status updates (QUEUED → BREWING → READY).
 * Also refetches when the browser tab regains focus.
 */
export function useOrders() {
  return useQuery({
    queryKey: queryKeys.orders,
    queryFn: fetchOrders,
    refetchInterval: 10_000,
  });
}

/**
 * Mutation to place a new order.
 * On success it automatically invalidates the orders cache
 * so the list refreshes without a manual refetch.
 */
export function useCreateOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (menuItemName: string) => createOrder(menuItemName),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.orders });
    },
  });
}

/**
 * Fetches a single order by ID.
 */
export function useOrder(orderId: string) {
  return useQuery({
    queryKey: queryKeys.order(orderId),
    queryFn: () => getOrder(orderId),
    enabled: !!orderId,
  });
}

/**
 * Mutation to update an order's status.
 * Invalidates both the orders list and the specific order cache.
 */
export function useUpdateOrderStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ orderId, status }: { orderId: string; status: string }) =>
      updateOrderStatus(orderId, status),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.orders });
      queryClient.invalidateQueries({
        queryKey: queryKeys.order(variables.orderId),
      });
    },
  });
}

/**
 * Mutation to delete an order.
 * Invalidates the orders list cache on success.
 */
export function useDeleteOrder() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (orderId: string) => deleteOrder(orderId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.orders });
    },
  });
}
