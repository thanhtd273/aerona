package com.thanhtd.aerona.base.service;

import java.util.List;

public interface RedisService {

    String createKey(String prefix, String key);
    void saveList(String key, List<Object> objects, int secondExpires);

    void saveObject(String key, Object object, int secondExpires);

    void save(String key, String value, int secondExpires);

    String getValue(String key);
    <T> T getObject(String key, Class<T> tClass);

    <T> List<T> getListObjects(String key, Class<T> tClass);

    void delete(String key);

}
