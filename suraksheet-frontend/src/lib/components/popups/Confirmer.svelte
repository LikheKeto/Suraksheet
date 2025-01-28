<script>
    import { Button, Modal } from 'flowbite-svelte';
    import { ExclamationCircleOutline } from 'flowbite-svelte-icons';
    export let message = ""
    export let showConfirmer = false;
    export let onConfirmation;
  </script>

  <Modal bind:open={showConfirmer} size="xs" autoclose>
    <div class="text-center">
      <ExclamationCircleOutline class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
      <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">{message}</h3>
      <Button on:click={onConfirmation} color="red" class="me-2">Yes, I'm sure</Button>
      <Button color="alternative">No, cancel</Button>
    </div>
  </Modal>