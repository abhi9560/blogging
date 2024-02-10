# Simple Blogging Platform API
This is a RESTful API for a simple blogging platform built using Go. It allows users to perform CRUD operations (Create, Read, Update, Delete) on blog posts.  <br />

Remember to replace "username:password@tcp(localhost:3306)/dbname" with your actual database connection details. Also, make sure to handle errors properly and add appropriate error messages in the responses.  <br />

# GET Request:
Method: GET  <br />
URL: http://localhost:8080/posts  <br />
Description: Retrieve all blog posts.  <br />
In Postman, simply select the request you created for retrieving all blog posts and click the "Send" button. Postman will send a GET request to http://localhost:8080/posts and display the response.  <br />

# POST Request:
Method: POST  <br />
URL: 'http://localhost:8080/posts'  <br />
Body (JSON):  <br />
  {  <br />
    "title": "New Blog Post",  <br />
    "content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit."  <br />
  }  <br />
Description: Create a new blog post.  <br />
In Postman, create a new request with the method set to POST and the URL set to 'http://localhost:8080/posts'. Add the above JSON body to the request body. Click the "Send" button to create a new blog post.  <br />

# PUT Request:
Method: PUT <br />
URL: http://localhost:8080/posts/{id} <br />
Body (JSON): <br />
  { <br />
    "title": "Updated Blog Post", <br />
    "content": "Updated content here." <br />
  } <br />
Description: Update an existing blog post. <br />
Replace {id} in the URL with the ID of the blog post you want to update. Set the method to PUT and enter the updated data in the request body as JSON. Click the "Send" button to update the blog post. <br />

# DELETE Request:
Method: DELETE <br />
URL: 'http://localhost:8080/posts/{id}' <br />
Description: Delete a blog post. <br />
Replace {id} in the URL with the ID of the blog post you want to delete. Set the method to DELETE and click the "Send" button to delete the blog post. <br />
