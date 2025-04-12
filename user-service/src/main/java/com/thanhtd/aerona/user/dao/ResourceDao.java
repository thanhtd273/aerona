package com.thanhtd.aerona.user.dao;

import com.thanhtd.aerona.user.model.Resource;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface ResourceDao extends JpaRepository<Resource, String> {

    @Query("SELECT u FROM Resource u WHERE u.status <> 0")
    List<Resource> findAll();

    @Query("SELECT u FROM Resource u WHERE u.status <> 0 AND u.resourceId = :resourceId")
    Resource findByResourceId(@Param(value = "resourceId") String resourceId);

    @Query("SELECT u FROM Resource u WHERE u.status <> 0 AND u.path = :path")
    Resource findByPath(@Param(value = "path") String path);
}
