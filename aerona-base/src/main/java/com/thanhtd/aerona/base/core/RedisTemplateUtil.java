package com.thanhtd.aerona.base.core;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.redis.core.*;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Set;
import java.util.concurrent.TimeUnit;

@Component
public class RedisTemplateUtil<T>{
    private RedisTemplate<String, T> redisTemplate;

    private ValueOperations<String, T> valueOperations;

    private ListOperations<String, T> listOperations;

    private SetOperations<String, T> setOperations;

    private HashOperations<String, String, T> hashOperations;

    @Autowired
    public RedisTemplateUtil(RedisTemplate<String, T> redisTemplate) {
        this.redisTemplate = redisTemplate;
        this.valueOperations = redisTemplate.opsForValue();
        this.listOperations = redisTemplate.opsForList();
        this.hashOperations = redisTemplate.opsForHash();
    }

    public boolean exist(String key) {
        return Boolean.TRUE.equals(redisTemplate.hasKey(key));
    }

    public void set(String key, T value) {
        valueOperations.set(key, value);
    }

    public void delete(String key) {
        redisTemplate.delete(key);
    }

    public T get(String key) {
        return valueOperations.get(key);
    }

    public void expire(String key, long timeout, TimeUnit unit) {
        redisTemplate.expire(key, timeout, unit);
    }

    public void addList(String key, T value) {
        listOperations.leftPush(key, value);
    }

    public List<T> getListMembers(String key) {
        return listOperations.range(key, 0, -1);
    }

    public Long getListSize(String key) {
        return listOperations.size(key);
    }

    public void addToSet(String key, T... values) {
        setOperations.add(key, values);
    }

    public Set<T> getSetMembers(String key) {
        return setOperations.members(key);
    }

    public void hset(String key, String field, T value) {
        hashOperations.put(key, field, value);
    }

    public T hget(String key, String field) {
        return hashOperations.get(key, field);
    }

    public void hdelete(String key, String field) {
        hashOperations.delete(key, field);
    }

    public String createKey(String prefix, String key) {
        return String.format("%s:%s", prefix, key);
    }

    public void setWithExpire(String key, T value, long timeout, TimeUnit unit) {
        this.set(key, value);
        this.expire(key, timeout, unit);
    }
}
