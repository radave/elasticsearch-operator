package k8shandler

const defaultKibanaIndexMode = "unique"

const esYmlTmpl = `
cluster:
  name: ${CLUSTER_NAME}

script:
  inline: true
  stored: true

node:
  name: ${DC_NAME}
  master: ${IS_MASTER}
  data: ${HAS_DATA}
  max_local_storage_nodes: 1

network:
  host: 0.0.0.0

cloud:
  kubernetes:
    service: ${SERVICE_DNS}
    namespace: ${NAMESPACE}

discovery.zen:
  hosts_provider: kubernetes
  minimum_master_nodes: ${NODE_QUORUM}

gateway:
  recover_after_nodes: ${NODE_QUORUM}
  expected_nodes: ${RECOVER_EXPECTED_NODES}
  recover_after_time: ${RECOVER_AFTER_TIME}

path:
  data:
  {{.PathData}}

`

const secureYmlTmpl = `
io.fabric8.elasticsearch.kibana.mapping.app: /usr/share/elasticsearch/index_patterns/com.redhat.viaq-openshift.index-pattern.json
io.fabric8.elasticsearch.kibana.mapping.ops: /usr/share/elasticsearch/index_patterns/com.redhat.viaq-openshift.index-pattern.json
io.fabric8.elasticsearch.kibana.mapping.empty: /usr/share/elasticsearch/index_patterns/com.redhat.viaq-openshift.index-pattern.json

openshift.config:
  use_common_data_model: true
  project_index_prefix: "project"
  time_field_name: "@timestamp"

openshift.searchguard:
  keystore.path: /etc/elasticsearch/secret/admin.jks
  truststore.path: /etc/elasticsearch/secret/searchguard.truststore

openshift.operations.allow_cluster_reader: {{.AllowClusterReader}}

openshift.kibana.index.mode: {{.KibanaIndexMode}}

searchguard:
  authcz.admin_dn:
  - CN=system.admin,OU=OpenShift,O=Logging
  config_index_name: ".searchguard.${DC_NAME}"
  ssl:
    transport:
      enabled: true
      enforce_hostname_verification: false
      keystore_type: JKS
      keystore_filepath: /etc/elasticsearch/secret/searchguard.key
      keystore_password: kspass
      truststore_type: JKS
      truststore_filepath: /etc/elasticsearch/secret/searchguard.truststore
      truststore_password: tspass
    http:
      enabled: true
      keystore_type: JKS
      keystore_filepath: /etc/elasticsearch/secret/key
      keystore_password: kspass
      clientauth_mode: OPTIONAL
      truststore_type: JKS
      truststore_filepath: /etc/elasticsearch/secret/truststore
      truststore_password: tspass`

const defaultRootLogger = "console"

const log4j2PropertiesTmpl = `
status = error

# log action execution errors for easier debugging
logger.action.name = org.elasticsearch.action
logger.action.level = debug

appender.console.type = Console
appender.console.name = console
appender.console.layout.type = PatternLayout
appender.console.layout.pattern = [%d{ISO8601}][%-5p][%-25c{1.}] %marker%m%n

appender.rolling.type = RollingFile
appender.rolling.name = rolling
appender.rolling.fileName = ${sys:es.logs.base_path}${sys:file.separator}${sys:es.logs.cluster_name}.log
appender.rolling.layout.type = PatternLayout
appender.rolling.layout.pattern = [%d{ISO8601}][%-5p][%-25c{1.}] %marker%.-10000m%n
appender.rolling.filePattern = ${sys:es.logs.base_path}${sys:file.separator}${sys:es.logs.cluster_name}-%d{yyyy-MM-dd}.log
appender.rolling.policies.type = Policies
appender.rolling.policies.time.type = TimeBasedTriggeringPolicy
appender.rolling.policies.time.interval = 1
appender.rolling.policies.time.modulate = true

rootLogger.level = info
rootLogger.appenderRef.{{.RootLogger}}.ref = {{.RootLogger}}

appender.deprecation_rolling.type = RollingFile
appender.deprecation_rolling.name = deprecation_rolling
appender.deprecation_rolling.fileName = ${sys:es.logs.base_path}${sys:file.separator}${sys:es.logs.cluster_name}_deprecation.log
appender.deprecation_rolling.layout.type = PatternLayout
appender.deprecation_rolling.layout.pattern = [%d{ISO8601}][%-5p][%-25c{1.}] %marker%.-10000m%n
appender.deprecation_rolling.filePattern = ${sys:es.logs.base_path}${sys:file.separator}${sys:es.logs.cluster_name}_deprecation-%i.log.gz
appender.deprecation_rolling.policies.type = Policies
appender.deprecation_rolling.policies.size.type = SizeBasedTriggeringPolicy
appender.deprecation_rolling.policies.size.size = 1GB
appender.deprecation_rolling.strategy.type = DefaultRolloverStrategy
appender.deprecation_rolling.strategy.max = 4

logger.deprecation.name = org.elasticsearch.deprecation
logger.deprecation.level = warn
logger.deprecation.appenderRef.deprecation_rolling.ref = deprecation_rolling
logger.deprecation.additivity = false

appender.index_search_slowlog_rolling.type = RollingFile
appender.index_search_slowlog_rolling.name = index_search_slowlog_rolling
appender.index_search_slowlog_rolling.fileName = ${sys:es.logs.base_path}${sys:file.separator}${sys:es.logs.cluster_name}_index_search_slowlog.log
appender.index_search_slowlog_rolling.layout.type = PatternLayout
appender.index_search_slowlog_rolling.layout.pattern = [%d{ISO8601}][%-5p][%-25c] %marker%.-10000m%n
appender.index_search_slowlog_rolling.filePattern = ${sys:es.logs.base_path}${sys:file.separator}${sys:es.logs.cluster_name}_index_search_slowlog-%d{yyyy-MM-dd}.log
appender.index_search_slowlog_rolling.policies.type = Policies
appender.index_search_slowlog_rolling.policies.time.type = TimeBasedTriggeringPolicy
appender.index_search_slowlog_rolling.policies.time.interval = 1
appender.index_search_slowlog_rolling.policies.time.modulate = true

logger.index_search_slowlog_rolling.name = index.search.slowlog
logger.index_search_slowlog_rolling.level = trace
logger.index_search_slowlog_rolling.appenderRef.index_search_slowlog_rolling.ref = index_search_slowlog_rolling
logger.index_search_slowlog_rolling.additivity = false

appender.index_indexing_slowlog_rolling.type = RollingFile
appender.index_indexing_slowlog_rolling.name = index_indexing_slowlog_rolling
appender.index_indexing_slowlog_rolling.fileName = ${sys:es.logs.base_path}${sys:file.separator}${sys:es.logs.cluster_name}_index_indexing_slowlog.log
appender.index_indexing_slowlog_rolling.layout.type = PatternLayout
appender.index_indexing_slowlog_rolling.layout.pattern = [%d{ISO8601}][%-5p][%-25c] %marker%.-10000m%n
appender.index_indexing_slowlog_rolling.filePattern = ${sys:es.logs.base_path}${sys:file.separator}${sys:es.logs.cluster_name}_index_indexing_slowlog-%d{yyyy-MM-dd}.log
appender.index_indexing_slowlog_rolling.policies.type = Policies
appender.index_indexing_slowlog_rolling.policies.time.type = TimeBasedTriggeringPolicy
appender.index_indexing_slowlog_rolling.policies.time.interval = 1
appender.index_indexing_slowlog_rolling.policies.time.modulate = true

logger.index_indexing_slowlog.name = index.indexing.slowlog.index
logger.index_indexing_slowlog.level = trace
logger.index_indexing_slowlog.appenderRef.index_indexing_slowlog_rolling.ref = index_indexing_slowlog_rolling
logger.index_indexing_slowlog.additivity = false`
