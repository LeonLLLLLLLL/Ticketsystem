<script lang="ts">
  import { fade } from 'svelte/transition';
  import { onDestroy } from 'svelte';
  import { toastStore, type ToastMessage } from '../stores/toastStore';

  let messages: ToastMessage[] = [];

  const unsubscribe = toastStore.subscribe((val) => {
    messages = val;
  });

  onDestroy(() => {
    unsubscribe();
  });

  function remove(id: number) {
    toastStore.remove(id);
  }
</script>

<style>
  .toast-container {
    position: fixed;
    top: 1rem;
    right: 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    z-index: 1000;
  }

  .toast {
    background-color: #333;
    color: white;
    padding: 0.75rem 1.25rem;
    border-radius: 5px;
    cursor: pointer;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
    max-width: 300px;
    font-size: 0.95rem;
  }

  .toast.success {
    background-color: #28a745;
  }

  .toast.error {
    background-color: #dc3545;
  }

  .toast.info {
    background-color: #17a2b8;
  }
</style>

<div class="toast-container">
  {#each messages as msg (msg.id)}
    <div
      class="toast {msg.type}"
      on:click={() => remove(msg.id)}
      in:fade
      out:fade
    >
      {msg.text}
    </div>
  {/each}
</div>
