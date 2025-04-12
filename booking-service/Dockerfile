FROM maven:3.9.6-eclipse-temurin-17 AS builder
WORKDIR /build

COPY pom.xml /build/pom.xml
COPY settings.xml /build/settings.xml
COPY src /build/src

RUN mvn clean install -DskipTests -s /build/settings.xml
RUN mkdir -p target/dependency && (cd target/dependency; jar -xf ../*.jar)

FROM eclipse-temurin:17-jre-alpine
VOLUME /tmp
ARG DEPENDENCY=/build/target/dependency

COPY --from=builder ${DEPENDENCY}/BOOT-INF/lib /app/lib
COPY --from=builder ${DEPENDENCY}/META-INF /app/META-INF
COPY --from=builder ${DEPENDENCY}/BOOT-INF/classes /app

VOLUME /tmptarget/*.jar
EXPOSE 8085
ENTRYPOINT ["java","-cp","app:app/lib/*", "com.thanhtd.aerona.booking.Application"]
