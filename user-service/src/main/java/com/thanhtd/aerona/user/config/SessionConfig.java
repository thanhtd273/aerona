package com.thanhtd.aerona.user.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.serializer.GenericJackson2JsonRedisSerializer;
import org.springframework.data.redis.serializer.RedisSerializer;
import org.springframework.session.MapSession;
import org.springframework.session.config.SessionRepositoryCustomizer;
import org.springframework.session.data.redis.RedisSessionMapper;
import org.springframework.session.data.redis.RedisSessionRepository;
import org.springframework.session.data.redis.config.annotation.web.http.EnableRedisHttpSession;

import java.util.Map;
import java.util.function.BiFunction;

@Slf4j
@Configuration
@EnableRedisHttpSession
public class SessionConfig {

    @Bean
    public RedisSerializer<Object> springSessionDefaultRedisSerializer() {
        ObjectMapper mapper = new ObjectMapper();
        return new GenericJackson2JsonRedisSerializer(mapper);
    }
    @Bean
    SessionRepositoryCustomizer<RedisSessionRepository> redisSessionRepositoryCustomizer() {
        return (redisSessionRepository) -> redisSessionRepository
                .setRedisSessionMapper(new SafeRedisSessionMapper(redisSessionRepository));
    }

    static class SafeRedisSessionMapper implements BiFunction<String, Map<String, Object>, MapSession> {

        private final RedisSessionMapper delegate = new RedisSessionMapper();

        private final RedisSessionRepository sessionRepository;

        SafeRedisSessionMapper(RedisSessionRepository sessionRepository) {
            this.sessionRepository = sessionRepository;
        }

        @Override
        public MapSession apply(String sessionId, Map<String, Object> map) {
            try {
                return this.delegate.apply(sessionId, map);
            }
            catch (IllegalStateException ex) {
                this.sessionRepository.deleteById(sessionId);
                return null;
            }
        }

    }
}
