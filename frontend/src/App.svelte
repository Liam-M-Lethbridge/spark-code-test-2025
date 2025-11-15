<script lang="ts">
  import Todo from "./lib/Todo.svelte";
  import type { TodoItem } from "./lib/types";
  let title = $state("");
  let description = $state("");
  let todos: TodoItem[] = $state([]);

  async function fetchTodos() {
    try {
      const response = await fetch("http://localhost:8080/");
      if (response.status !== 200) {
        console.error("Error fetching data. Response status not 200");
        return;
      }

      todos = await response.json();
    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  async function addTodo(event: SubmitEvent) {
    // in my opinon, submitting the form normally and causing a reload is janky and annoying. I think preventing the default is more seemless
    event.preventDefault()

    // make sure the inputs are valid
    if(title === "" || description === "" ){
      // ideally I would add something like a pop-up to let the user know but I didn't want to mess with the layout of the webpage
      console.log("Invalid inputs")
      return
    }

    try{
      const response = await fetch("http://localhost:8080/",{
      method: "POST",
      headers: { "Content-Type": "application/json"},
      body: JSON.stringify({title, description})
    });
      // If we get the item back, we know it has been processes correctly and can push to client list without needing to get the entire list back
      const todo = await response.json() as TodoItem
      // check if todo is null. Very annoying svelte reactive property
      if(todos){
        todos.push(todo)
      }else{
        todos = [todo]
      }  
      title = ""
      description = ""
  } catch (e) {
      console.error("Error processing Todo.")
  }
}

  // Initially fetch todos on page load
  $effect(() => {
    fetchTodos();
  });

/*
I was originally going to add functionality like "crossing a task off" 
but to have this task permanently crossed off I would need extra functionality on the backend   
and the OpenAPI spec does not include any functionality for marking a task as completed or removing it from the list so I left it
*/
</script>

<main class="app">
  <header class="app-header">
    <h1>TODO</h1>
  </header>

  <div class="todo-list">
    {#each todos as todo}
      <Todo title={todo.title} description={todo.description} />
    {/each}
  </div>

  <h2 class="todo-list-form-header">Add a Todo</h2>
  <form class="todo-list-form" onsubmit={addTodo}>
    <input placeholder="Title" name="title" bind:value={title} minlength="1"/>
    <input placeholder="Description" name="description" id="description" bind:value={description} minlength="1"/>
    <button>Add Todo</button>
  </form>
</main>

<style>
  .app {
    color: white;
    background-color: #282c34;

    text-align: center;
    font-size: 24px;

    min-height: 100vh;
    padding: 20px;
  }

  .app-header {
    font-size: calc(10px + 4vmin);
    margin-top: 50px;
  }

  .todo-list {
    margin: 50px 100px 0px 100px;
  }

  .todo-list-form-header {
    margin-top: 100px;
  }

  .todo-list-form {
    margin-top: 10px;
  }
</style>
