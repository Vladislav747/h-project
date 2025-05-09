INSERT INTO companies (id, name, description)
VALUES (:id, :name, :description)
    RETURNING id;