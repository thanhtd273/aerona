package com.thanhtd.aerona.base.service.impl;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.thanhtd.aerona.base.core.RedisTemplateUtil;
import com.thanhtd.aerona.base.service.RedisService;
import jodd.util.StringPool;
import jodd.util.StringUtil;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.util.ObjectUtils;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;

@Service
@RequiredArgsConstructor
public class RedisServiceImpl implements RedisService {
    private static final Logger logger = LoggerFactory.getLogger(RedisServiceImpl.class);

    private final RedisTemplateUtil<String> redisTemplateUtil;

    private final ObjectMapper jsonMapper;

    @Override
    public String createKey(String prefix, String key) {
        return redisTemplateUtil.createKey(prefix, key);
    }

    @Override
    public void saveList(String key, List<Object> objects, int secondExpires) {
        try {
            String value = StringPool.EMPTY;
            if (!ObjectUtils.isEmpty(objects)) {
                value = jsonMapper.writeValueAsString(objects);
            }
            save(key, value, secondExpires);
        } catch (JsonProcessingException e) {
            logger.error("Cannot save list of object to redis cache");
        }

    }

    @Override
    public void saveObject(String key, Object object, int secondExpires) {
        try {
            String value = StringPool.EMPTY;
            if (!ObjectUtils.isEmpty(object)) {
                value = jsonMapper.writeValueAsString(object);
            }
            save(key, value, secondExpires);
        } catch (JsonProcessingException e) {
            logger.error("Cannot save object to redis cache");
        }

    }

    @Override
    public void save(String key, String value, int secondExpires) {
        redisTemplateUtil.setWithExpire(key, value, secondExpires, TimeUnit.SECONDS);
    }

    @Override
    public String getValue(String key) {
        if (redisTemplateUtil.exist(key))
            return redisTemplateUtil.get(key);
        return null;
    }

    @Override
    public <T> T getObject(String key, Class<T> tClass) {
        try {
            String value = getValue(key);
            if (StringUtil.isBlank(value))
                return null;
            return jsonMapper.readValue(value, tClass);

        } catch (IOException e) {
            logger.error("Get object from redis fail, error: {}", e.getMessage());
            return null;
        }
    }

    @Override
    public <T> List<T> getListObjects(String key, Class<T> tClass) {
        try {
            String value = getValue(key);
            if (StringUtil.isBlank(value))
                return new ArrayList<>();
            return jsonMapper.readValue(value, jsonMapper.getTypeFactory().constructCollectionType(List.class, tClass));
        } catch (IOException e) {
            logger.error("Get list of object from redis cache fail, error: {}", e.getMessage());
            return new ArrayList<>();
        }
    }

    @Override
    public void delete(String key) {
        if (redisTemplateUtil.exist(key))
            redisTemplateUtil.delete(key);
    }
}
