# Coding exercise

If you are applying for a position at Sourcegraph, read [the instructions](INSTRUCTIONS.md). You
will add content to this README as you complete parts of the exercise.

If you are a Sourcegraph team member, send a zip of this directory to the candidate.
Prior to sending them the directory, ensure they have [Docker](https://docs.docker.com/get-docker/) and [docker-compose](https://docs.docker.com/compose/install/) installed.

(__Note:__ This is a separate binary, not `docker compose` that is bundled with the Docker for Desktop installation ) .


## CANDIDATE SUBMISSION:

To start, I did attempt to install and use Docker, but this is the first time I've touched Docker, so I was a little confused how to use it.

So I am able to say a conceptual view of the issue the code has.

1. Your draft answer to the customer (write it as if you are writing a real customer) AND a note to your internal teammates on how you got to the answer (write this as if you are writing your teammates to show them your thought process) AND how long the Docker install took (if relevant)

### My answer to the customer: 
"The application clearly does not allow for you to enter any input at all. My team and I would need to build the backend part of it, that allows the input to be stored and saved. So that way when you enter anything into the todo field, it can be shown on the page and still persist upon refresh. We will update you within this week as soon as we can with more updates.

Best,
Engineering team at Sourcegraph"

### My answer to my teammates:
"Based on customer feedback, the do list application has couple bugs:
- So the input field has a bug, input cannot be entered. There may be something wrong with the routing, or with the data that the input field is seeking.
- New to do list items need to persist on the page on refresh. Since this appears to be a simple DOM manupulation, we need to build a local storage to save that data.

For local storage, a localstorage would be triggered on every new todo list item entered, adding it to storage. It would stored i n local stroage at a list item element. It would be identified by the id, document.getElementById("add-task");, then added to localstorage.

Each time a new item is added to the list, there needs to be a submit and click action include, to tell the function what to do when that action happend. The page would automatically refresh, disrupting the local storage, so we need to add e.preventDefault(); to the submit function for the addEventListener function to the form. 

The updateToDO and addTodo functions are not really returning anything. It should use .innerText to update the DOM, so we can see the changes on the page. 

As a stretch goal, if time permits, we can also include a "remove" button, to remove an item from the list. This would simply be a .createElement('button') that would be appended onto a new entry. 