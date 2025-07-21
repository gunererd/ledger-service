// MongoDB initialization script
db = db.getSiblingDB('ledger');

// Create collections
db.createCollection('transactions');
db.createCollection('balances');

// Create indexes for better performance
db.transactions.createIndex({ "customer.id": 1 });
db.transactions.createIndex({ "restaurant.id": 1 });
db.transactions.createIndex({ "type": 1 });
db.balances.createIndex({ "userid": 1 }, { unique: true });

print('Database initialized successfully');