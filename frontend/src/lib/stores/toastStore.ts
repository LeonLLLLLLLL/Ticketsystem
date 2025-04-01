import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'info';

export interface ToastMessage {
  id: number;
  text: string;
  type: ToastType;
}

function createToastStore() {
  const { subscribe, update } = writable<ToastMessage[]>([]);
  let idCounter = 0;

  function push(text: string, type: ToastType = 'info', duration = 3000) {
    const id = ++idCounter;
    const newToast: ToastMessage = { id, text, type };

    update((toasts) => [...toasts, newToast]);

    setTimeout(() => {
      remove(id);
    }, duration);
  }

  function remove(id: number) {
    update((toasts) => toasts.filter((t) => t.id !== id));
  }

  return {
    subscribe,
    push,
    remove
  };
}

export const toastStore = createToastStore();
