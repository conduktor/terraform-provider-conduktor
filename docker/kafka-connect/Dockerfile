FROM confluentinc/cp-kafka-connect:7.7.1

# Add Apache Kafka connect-tools jar (contains MockSourceConnector used in tests).
# Confluent removed it from their distribution in CP 7.4+; the open-source jar is still published.
USER root
RUN curl -fsSL \
    "https://repo1.maven.org/maven2/org/apache/kafka/connect-test-plugins/3.7.1/connect-test-plugins-3.7.1.jar" \
    -o /usr/share/java/kafka/connect-test-plugins-3.7.1.jar
USER appuser
