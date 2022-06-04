db = db.getSiblingDB('homework');

db.createUser(
        {
            user: 'fredy',
            pwd: 'homework',
            roles: [
                {
                    role: 'readWrite',
                    db: 'homework'
                }
            ]
        }
);