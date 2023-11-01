## First problem (Solved ✅)

I want a structure that does the following things

1. Create a mapping for props
2. Able to check if the current prop is mapped or not
3. Able to access the actual value of the prop from the mapping

for mapping I need a uniquely identifiable key but also that key should be
recreated if i need
(can't use value because that can change, can't use key because there can be duplicates)
can use pointers here
This assumption was wrong because the key or prop name by nature is unique for a component so a simple mapping with some abstraction solved the problem it is written below.

Pointer approach:
instead of having the key as a string, we can make it a pointer of a string or any entity

Ended up using a map that would map all the dynamic props with a key to the actual value and it worked for me.
Also since Golang has no getter thing, I implemented an abstraction for the getter that basically just uses `Get` as a function and whenever we use it, it makes a proxy or a mapping, here the point to note is that whenever we read something through `Get` it will be a dynamic prop.

## Second problem or Chain of thought

I want a compiler/transpiler something like Babel. But before that, I also need to decide what type of syntax and semantics I want for the vdom itself.

Here are a few options:

1. JSX style code similar to react:

   ```jsx
   const Button = ({ dynamicContent }) => {
     return <button>{dynamicContent}</button>;
   };
   ```

2. Svelte style code

   ```svelte
   <style>
       /* Some css here */
   </style>
   <script>
       // some javascript
       let dynamicContent = "something"
   </script>

   <button>
       {dynamicContent}
   </button>

   ```

   In our case the script might be replaced with some other tag or maybe have a attribute to support go code.

   ```html
   <style>
     /* Some css here */
   </style>
   <script type="go/script">
     // some javascript
     dynamicContent := "something"
   </script>

   <button>{dynamicContent}</button>
   ```

3. Template style (like jinja or any basic templating engine)
   ```go
   func Button(dynamicContent string) Component{
      return `<button>{{dynamicContent}}</button>`
   }
   ```
