## Flow


All API calls in the frontend are going through React Query. Here's the full picture:

### Layer 1: api.ts — Pure fetch functions (6 total)

| Function                | Connect RPC Endpoint               | Used in Hook?               |
|-------------------------|----------------------------------|----------------------------|
| fetchMenu()             | menu.MenuService/GetMenu          | ✅ useMenu                 |
| fetchOrders()           | brew.BrewService/ListOrders       | ✅ useOrders               |
| createOrder()           | brew.BrewService/OrderDrink       | ✅ useCreateOrder          |
| getOrder()              | brew.BrewService/GetOrder         | ✅ useOrder                |
| updateOrderStatus()     | brew.BrewService/UpdateOrderStatus | ✅ useUpdateOrderStatus    |
| deleteOrder()           | brew.BrewService/DeleteOrder      | ✅ useDeleteOrder          |

### Layer 2: hooks.ts — React Query hooks (6 total)

| Hook                        | Type         | Used in Component?              |
|-----------------------------|-------------|--------------------------------|
| useMenu()                   | useQuery    | ✅ Menu.tsx                    |
| useOrders()                 | useQuery    | ✅ Orders.tsx                  |
| useCreateOrder()            | useMutation | ✅ OrderForm.tsx               |
| useOrder()                  | useQuery    | Available (not used yet)       |
| useUpdateOrderStatus()      | useMutation | Available (not used yet)       |
| useDeleteOrder()            | useMutation | Available (not used yet)       |

Layer 3: Components — Only import from hooks, never from api.ts directly
✅ No component imports from api.ts directly — all go through React Query hooks. Clean separation.

The three new hooks (useOrder, useUpdateOrderStatus, useDeleteOrder) are ready to use whenever you add UI for those features (e.g., a status dropdown or delete button on each order row).

# React + TypeScript + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Babel](https://babeljs.io/) (or [oxc](https://oxc.rs) when used in [rolldown-vite](https://vite.dev/guide/rolldown)) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## React Compiler

The React Compiler is not enabled on this template because of its impact on dev & build performances. To add it, see [this documentation](https://react.dev/learn/react-compiler/installation).

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type-aware lint rules:

```js
export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...

      // Remove tseslint.configs.recommended and replace with this
      tseslint.configs.recommendedTypeChecked,
      // Alternatively, use this for stricter rules
      tseslint.configs.strictTypeChecked,
      // Optionally, add this for stylistic rules
      tseslint.configs.stylisticTypeChecked,

      // Other configs...
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```

You can also install [eslint-plugin-react-x](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://github.com/Rel1cx/eslint-react/tree/main/packages/plugins/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])
```
