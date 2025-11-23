const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args));

async function testRegister() {
  const query = `
    mutation {
      register(input: {
        email: "newuser@example.com"
        password: "password123"
        name: "New User"
      }) {
        token
        user {
          id
          email
          name
        }
      }
    }
  `;

  try {
    const response = await fetch('http://localhost:8080/query', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        query,
      }),
    }).then(res => res.json());

    console.log(JSON.stringify(response, null, 2));
  } catch (error) {
    console.error('Error:', error);
  }
}

testRegister();