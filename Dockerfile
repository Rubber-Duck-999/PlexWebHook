FROM arm32v7/maven:latest AS build
WORKDIR /app
COPY . .
RUN mvn clean package -DskipTests

FROM arm32v7/openjdk:latest
WORKDIR /app
COPY --from=build /app/target/plexwebhook-1.0.0.jar app.jar
EXPOSE 8000
ENTRYPOINT ["java", "-jar", "app.jar"]