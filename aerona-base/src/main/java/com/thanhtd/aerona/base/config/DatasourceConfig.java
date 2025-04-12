package com.thanhtd.aerona.base.config;

import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import javax.sql.DataSource;
import java.util.Properties;

@Configuration
public class DatasourceConfig {

    @Value("${spring.datasource.username}")
    private String username;

    @Value("${spring.data.redis.password}")
    private String password;

    @Value("${spring.datasource.url}")
    private String url;

    @Value("${spring.datasource.driver-class-name}")
    private String driverClassName;

    @Value("${spring.sql.init.continue-on-error}")
    private String continueOnError;

    @Value("${spring.datasource.hikari.minimum-idle}")
    private int minimumIdle;

    @Value("${spring.datasource.hikari.maximum-pool-size}")
    private int maximumPoolSize;

    @Value("${spring.datasource.hikari.connection-timeout}")
    private long connectionTimeout;

    @Value("${spring.datasource.hikari.idle-timeout}")
    private long idleTimeout;

    @Value("${spring.datasource.hikari.max-lifetime}")
    private long maxLifetime;

    @Bean
    public DataSource primaryDataSource() {
        Properties datasourceProps = new Properties();
        datasourceProps.setProperty("user", username);
        datasourceProps.setProperty("password", password);
        datasourceProps.setProperty("url", url);
        datasourceProps.setProperty("continue-on-error", continueOnError);

        Properties configProps = new Properties();
        configProps.setProperty("driverClassName", driverClassName);
        configProps.setProperty("jdbcUrl", url);

        HikariConfig hikariConfig = new HikariConfig(configProps);
        hikariConfig.setDataSourceProperties(datasourceProps);
        hikariConfig.setMinimumIdle(minimumIdle);
        hikariConfig.setMaximumPoolSize(maximumPoolSize);
//        hikariConfig.setConnectionTimeout(connectionTimeout);
//        hikariConfig.setIdleTimeout(idleTimeout);
        hikariConfig.setMaxLifetime(maxLifetime);
        return new HikariDataSource(hikariConfig);
    }
}
