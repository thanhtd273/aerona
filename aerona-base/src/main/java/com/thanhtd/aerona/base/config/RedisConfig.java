package com.thanhtd.aerona.base.config;

import lombok.RequiredArgsConstructor;
import org.redisson.Redisson;
import org.redisson.api.RedissonClient;
import org.redisson.spring.data.connection.RedissonConnectionFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.AutoConfigureBefore;
import org.springframework.boot.autoconfigure.condition.ConditionalOnClass;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.boot.autoconfigure.data.redis.RedisProperties;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.connection.*;
import org.springframework.data.redis.connection.jedis.JedisClientConfiguration;
import org.springframework.data.redis.connection.jedis.JedisConnectionFactory;
import org.springframework.data.redis.core.RedisOperations;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.serializer.*;
import redis.clients.jedis.JedisPoolConfig;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

@Configuration
@ConditionalOnClass({Redisson.class, RedisOperations.class})
@AutoConfigureBefore({RedisConfiguration.class})
@EnableConfigurationProperties({RedisProperties.class})
@RequiredArgsConstructor
public class RedisConfig {

    private static final Logger logger = LoggerFactory.getLogger(RedisConfig.class);

    /** The Constant Redis prefix. */
    private static final String REDIS_PREFIX = "redis://";

    /** The Constant Redis SSL prefix. */
    private static final String REDIS_SSL_PREFIX = "rediss://";

    @Value("${spring.data.redis.host}")
    private String redisHost;

    @Value("${spring.data.redis.port}")
    private int redisPort;

    @Value("${spring.data.redis.password}")
    private String redisPassword;

    @Value("${redis.connection.timeout}")
    private int redisConnectionTimeout;

    @Value("${redis.read.timeout}")
    private int redisReadTimeout;

    @Value("${redis.max.wait.millis}")
    private int maxWaitMillis;

    @Value("${redis.max.total}")
    private int redisMaxTotal;

    @Value("${redis.min.idle}")
    private int redisMinIdle;

    @Value("${redis.max.idle}")
    private int redisMaxIdle;

    @Value("${spring.data.redis.client-type}")
    private String clientType;

    @Value("${spring.data.redis.cluster.max-redirects}")
    private int clusterMaxDirects;

    private final RedisProperties redisProperties;

    @Bean
    public RedisConnectionFactory redisConnectionFactory() {

        logger.info("Create connection factory with stand alone configuration");
        RedisStandaloneConfiguration redisStandaloneConfiguration = new RedisStandaloneConfiguration(redisHost, redisPort);
        redisStandaloneConfiguration.setPassword(RedisPassword.of(redisPassword));
        redisStandaloneConfiguration.setDatabase(redisProperties.getDatabase());
        JedisClientConfiguration jedisClientConfig = getJedisClientConfig();
        JedisConnectionFactory jedisConnectionFactory = new JedisConnectionFactory(redisStandaloneConfiguration, jedisClientConfig);
        jedisConnectionFactory.afterPropertiesSet();
        return jedisConnectionFactory;
    }

    private JedisClientConfiguration getJedisClientConfig(){
        return JedisClientConfiguration
                .builder()
                .connectTimeout(Duration.ofSeconds(redisConnectionTimeout))
                .readTimeout(Duration.ofSeconds(redisReadTimeout))
                .usePooling()
                .poolConfig(jedisPoolConfig())
                .build();
    }

    @Bean
    public JedisPoolConfig jedisPoolConfig() {
        JedisPoolConfig jedisPoolConfig = new JedisPoolConfig();
        jedisPoolConfig.setMaxWait(Duration.ofSeconds(maxWaitMillis));
        jedisPoolConfig.setMaxTotal(redisMaxTotal);
        jedisPoolConfig.setMinIdle(redisMinIdle);
        jedisPoolConfig.setMaxIdle(redisMaxIdle);
        jedisPoolConfig.setTestOnBorrow(true);
        jedisPoolConfig.setTestOnCreate(true);
        jedisPoolConfig.setTestWhileIdle(true);
        jedisPoolConfig.setMinEvictableIdleDuration(Duration.ofSeconds(1800000));
        jedisPoolConfig.setTimeBetweenEvictionRuns(Duration.ofSeconds(30));
        jedisPoolConfig.setNumTestsPerEvictionRun(3);
        jedisPoolConfig.setBlockWhenExhausted(true);

        return jedisPoolConfig;
    }

    @Bean
    public RedisTemplate<Object, Object> redisTemplate() {
        final RedisTemplate<Object, Object> redisTemplate = new RedisTemplate<>();
        redisTemplate.setKeySerializer(new StringRedisSerializer());
        redisTemplate.setValueSerializer(new Jackson2JsonRedisSerializer<>(Object.class));
        redisTemplate.setHashKeySerializer(new GenericToStringSerializer<>(Object.class));
        redisTemplate.setHashValueSerializer(new GenericJackson2JsonRedisSerializer());
        redisTemplate.setConnectionFactory(redisConnectionFactory());
        return redisTemplate;
    }

    @Bean
    @ConditionalOnMissingBean({RedisConnectionFactory.class})
    public RedissonConnectionFactory redissonConnectionFactory(RedissonClient redisson) {
        return new RedissonConnectionFactory(redisson);
    }

//    @Bean(destroyMethod = "shutdown")
//    @ConditionalOnMissingBean({RedissonClient.class})
//    public RedissonClient redisson() {
//        Config config = new Config();
//        if (!redisProperties.getCluster().getNodes().isEmpty()) {
//            String[] listNodes = convert(redisProperties.getCluster().getNodes());
//            config.useClusterServers()
//                    .addNodeAddress(listNodes)
//                    .setConnectTimeout(redisConnectionTimeout)
//                    .setPassword(redisProperties.getPassword());
//        } else {
//            String prefix = this.redisProperties.getSsl().isEnabled() ? REDIS_SSL_PREFIX : REDIS_PREFIX;
//            config.useSingleServer()
//                    .setAddress(prefix + redisProperties.getHost() + ":" + redisProperties.getPort())
//                    .setConnectTimeout(redisConnectionTimeout)
//                    .setDatabase(redisProperties.getDatabase())
//                    .setPassword(redisProperties.getPassword());
//        }
//        return Redisson.create(config);
//    }

    private String[] convert(List<String> nodesObject) {
        List<String> listNode = new ArrayList<>(nodesObject.size());
        Iterator var3 = nodesObject.iterator();
        while (true) {
            while (var3.hasNext()) {
                String node = (String) var3.next();
                if (!node.startsWith(REDIS_PREFIX) && !node.startsWith(REDIS_SSL_PREFIX)) {
                    listNode.add(REDIS_PREFIX + node);
                } else {
                    listNode.add(node);
                }
            }
            return listNode.toArray(new String[listNode.size()]);
        }
    }
}
