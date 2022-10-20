DO $$
BEGIN

   IF NOT EXISTS
   	   (SELECT * FROM pg_tables
   		WHERE  schemaname = 'public'
   		AND    tablename  = 'location') THEN

        CREATE TABLE "location" (
                            "id" SERIAL PRIMARY KEY,
                            "username" varchar(16) UNIQUE NOT NULL,
                            "latitude" float8 NOT NULL,
                            "longitude" float8 NOT NULL,
                            "updated_at" timestamp NOT NULL
        );
    END IF;

   IF NOT EXISTS
   	   (SELECT * FROM pg_tables
   		WHERE  schemaname = 'public'
   		AND    tablename  = 'location_history') THEN

        CREATE TABLE "location_history" (
                                    "id" SERIAL PRIMARY KEY,
                                    "username" varchar(16) NOT NULL,
                                    "latitude" float8 NOT NULL,
                                    "longitude" float8 NOT NULL,
                                    "updated_at" timestamp NOT NULL,
                                    "distance" float8
        );
    END IF;

END;
$$;