const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args));

async function testGraphQL() {
  const query = `
    query {
      products {
        id
        name
        basePrice
      }
    }
  `;

  try {
    const response = await fetch('http://localhost:8081/query', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        query,
      }),
    }).then(res => res.json());

    console.log("GraphQL query result:");
    console.log(JSON.stringify(response, null, 2));
  } catch (error) {
    console.error('Error:', error);
  }
}

testGraphQL();

async function testCreateProductVariant() {
  // First, let's create a product
  const createProductQuery = `
    mutation {
      createProduct(input: {
        name: "Test T-Shirt"
        description: "A beautiful test t-shirt"
        designImageURL: "https://example.com/image.jpg"
        basePrice: 19.99
      }) {
        id
        name
        basePrice
      }
    }
  `;

  try {
    // Create product
    const productResponse = await fetch('http://localhost:8081/query', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        query: createProductQuery,
      }),
    }).then(res => res.json());

    console.log("Product creation result:");
    console.log(JSON.stringify(productResponse, null, 2));

    // If product creation was successful, create a variant
    if (productResponse.data && productResponse.data.createProduct) {
      const productId = productResponse.data.createProduct.id;
      
      const createVariantQuery = `
        mutation {
          createProductVariant(input: {
            productID: "${productId}"
            size: "M"
            color: "Red"
            priceModifier: 0.0
            sku: "TS-M-RED"
            stockQuantity: 100
          }) {
            id
            size
            color
            sku
            inventory {
              stockQuantity
            }
          }
        }
      `;

      const variantResponse = await fetch('http://localhost:8081/query', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          query: createVariantQuery,
        }),
      }).then(res => res.json());

      console.log("\nVariant creation result:");
      console.log(JSON.stringify(variantResponse, null, 2));
    }
  } catch (error) {
    console.error('Error:', error);
  }
}

testCreateProductVariant();

async function testInvalidToken() {
  const query = `
    query {
      products {
        id
        name
      }
    }
  `;

  try {
    const response = await fetch('http://localhost:8081/query', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer invalid-token'
      },
      body: JSON.stringify({
        query,
      }),
    });

    const result = await response.json();
    console.log("Response with invalid token:");
    console.log(JSON.stringify(result, null, 2));
  } catch (error) {
    console.error('Error:', error);
  }
}

testInvalidToken();