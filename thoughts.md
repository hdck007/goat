i want a structure that does the following things

1. Create a mapping for props
2. Able to check it the current prop is mapped or not
3. Able to access the actual value of the prop from the mapping

for mapping i need a uniquely identifiable key but also that key should be
recreated if i need
(cant use value because that can change, cant use key because there can be duplicates)
can use pointers here

Pointer approach:
instead of having the key as string we can make it a pointer of string or any entity
