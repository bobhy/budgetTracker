<script lang="ts">
   import { Button } from "$lib/components/ui/button";
   import { CleanDatabase } from "../../../wailsjs/go/main/App";

   let message = $state("");

   async function cleanDB() {
       try {
           await CleanDatabase();
           message = "Database cleaned successfully.";
       } catch (e) {
           message = "Error cleaning: " + String(e);
       }
   }
</script>

<div class="p-6">
   <h2 class="text-3xl font-bold mb-6">Database Administration</h2>
   
   <p class="mb-4 text-muted-foreground">
       Warning: Cleaning the database will remove all data and reset the schema.
   </p>

   <div class="flex gap-4 items-center">
       <Button variant="destructive" onclick={cleanDB}>Clean Database</Button>
       {#if message}
           <span class="text-sm border p-2 rounded bg-muted">{message}</span>
       {/if}
   </div>
</div>
