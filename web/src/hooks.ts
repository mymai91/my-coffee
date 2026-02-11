import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchMenu, fetchOrders, createOrder } from "./api";

// Query key constants — avoids typos and makes invalidation easy
export const queryKeys = {
  menu: ["menu"] as const,
  orders: ["orders"] as const,
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
