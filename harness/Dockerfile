FROM cassandra:3.0.18

# Add in the Jolokia agent
ADD https://repo1.maven.org/maven2/org/jolokia/jolokia-jvm/1.3.7/jolokia-jvm-1.3.7-agent.jar /usr/share/jolokia-jvm-1.3.7-agent.jar
RUN chmod 655 /usr/share/jolokia-jvm-1.3.7-agent.jar

ENV JVM_OPTS="-javaagent:/usr/share/jolokia-jvm-1.3.7-agent.jar=port=8778,host=0.0.0.0"

# 7000: intra-node communication
# 7001: TLS intra-node communication
# 7199: JMX
# 8778: Jolokia
# 9042: CQL
# 9160: thrift service
EXPOSE 7000 7001 7199 8778 9042 9160
CMD ["cassandra", "-f"]
