-- MySQL dump 10.14  Distrib 5.5.60-MariaDB, for Linux (x86_64)
--
-- Host: 10.57.242.40    Database: data4test
-- ------------------------------------------------------
-- Server version	5.7.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `api_definition`
--

CREATE DATABASE IF NOT EXISTS data4test;
USE data4test;

DROP TABLE IF EXISTS `api_definition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_definition` (
                                  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                  `api_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '接口ID',
                                  `api_module` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属模块',
                                  `api_desc` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '接口描述',
                                  `http_method` enum('get','post','put','delete') COLLATE utf8_unicode_ci DEFAULT 'get' COMMENT '请求方法',
                                  `path` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '请求路径',
                                  `header` longtext COLLATE utf8_unicode_ci COMMENT 'Header参数',
                                  `path_variable` longtext COLLATE utf8_unicode_ci COMMENT 'PATH参数',
                                  `query_parameter` longtext COLLATE utf8_unicode_ci COMMENT 'Query参数',
                                  `body` longtext COLLATE utf8_unicode_ci COMMENT 'Body参数',
                                  `response` longtext COLLATE utf8_unicode_ci COMMENT 'Resp参数',
                                  `version` int(11) DEFAULT NULL COMMENT '接口版本',
                                  `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                  `api_status` int(11) DEFAULT NULL COMMENT '接口状态，1:新增，2:被删除，3:被修改，4:保持原样',
                                  `change_content` longtext COLLATE utf8_unicode_ci COMMENT '变更内容',
                                  `check` varchar(4) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '规范检查',
                                  `api_check_fail_reason` longtext COLLATE utf8_unicode_ci COMMENT '规范检查失败原因',
                                  `auto` varchar(4) COLLATE utf8_unicode_ci DEFAULT '是' COMMENT '是否自动化',
                                  `remark` longtext COLLATE utf8_unicode_ci COMMENT '备注',
                                  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='接口定义';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `api_fuzzing_data`
--

DROP TABLE IF EXISTS `api_fuzzing_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_fuzzing_data` (
                                    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                    `data_desc` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '数据描述',
                                    `api_desc` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '接口描述',
                                    `api_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '关联接口',
                                    `api_module` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属模块',
                                    `header` longtext COLLATE utf8_unicode_ci COMMENT 'Header',
                                    `url_query` longtext COLLATE utf8_unicode_ci COMMENT 'UrlQuery',
                                    `body` longtext COLLATE utf8_unicode_ci COMMENT 'Body',
                                    `run_num` int(11) DEFAULT '1' COMMENT '执行次数',
                                    `expected_result` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '预期结果',
                                    `actual_result` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '实际结果',
                                    `result` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '测试结果',
                                    `fail_reason` longtext COLLATE utf8_unicode_ci COMMENT '失败原因',
                                    `response` longtext COLLATE utf8_unicode_ci COMMENT 'Response',
                                    `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='模糊数据';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `api_id_count`
--

DROP TABLE IF EXISTS `api_id_count`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_id_count` (
                                `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                `api_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '接口ID',
                                `api_desc` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '接口描述',
                                `run_times` int(11) NOT NULL COMMENT '执行次数',
                                `test_times` int(11) NOT NULL COMMENT '测试次数',
                                `pass_times` int(11) NOT NULL COMMENT '通过次数',
                                `fail_times` int(11) NOT NULL COMMENT '失败次数',
                                `untest_times` int(11) NOT NULL COMMENT '未测试次数',
                                `test_result` char(8) COLLATE utf8_unicode_ci NOT NULL COMMENT '测试结果',
                                `fail_reason` longtext COLLATE utf8_unicode_ci COMMENT '失败原因',
                                `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='接口统计';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `api_relation`
--

DROP TABLE IF EXISTS `api_relation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_relation` (
                                `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                `api_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '接口ID',
                                `api_desc` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '接口描述',
                                `api_module` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属模块',
                                `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                `auto` enum('yes','no') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'yes' COMMENT '是否自动化',
                                `pre_apis` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '前置接口',
                                `out_vars` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '提供变量关系',
                                `check_vars` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '校验变量转换关系',
                                `param_apis` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '依赖参数关联接口',
                                `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `api_id_app` (`api_id`,`app`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='接口关系';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `api_test_data`
--

DROP TABLE IF EXISTS `api_test_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_test_data` (
                                 `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                 `data_desc` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '数据描述',
                                 `api_desc` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '接口描述',
                                 `api_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '关联接口',
                                 `api_module` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属模块',
                                 `header` longtext COLLATE utf8_unicode_ci COMMENT 'Header',
                                 `url_query` longtext COLLATE utf8_unicode_ci COMMENT 'UrlQuery',
                                 `body` longtext COLLATE utf8_unicode_ci COMMENT 'Body',
                                 `run_num` int(11) DEFAULT '1' COMMENT '执行次数',
                                 `expected_result` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '预期结果',
                                 `actual_result` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '实际结果',
                                 `result` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '测试结果',
                                 `fail_reason` longtext COLLATE utf8_unicode_ci COMMENT '失败原因',
                                 `response` longtext COLLATE utf8_unicode_ci COMMENT 'Response',
                                 `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                 `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='测试数据';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `api_test_detail`
--

DROP TABLE IF EXISTS `api_test_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_test_detail` (
                                   `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                   `api_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '接口ID',
                                   `api_desc` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '接口描述',
                                   `data_desc` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
                                   `header` longtext COLLATE utf8_unicode_ci COMMENT 'Header',
                                   `url` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'URL',
                                   `body` longtext COLLATE utf8_unicode_ci COMMENT 'Body',
                                   `response` longtext COLLATE utf8_unicode_ci COMMENT 'Response',
                                   `fail_reason` longtext COLLATE utf8_unicode_ci COMMENT '失败原因',
                                   `test_result` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '测试结果',
                                   `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                   `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='测试详情';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `api_test_result`
--

DROP TABLE IF EXISTS `api_test_result`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `api_test_result` (
                                   `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                   `api_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '接口ID',
                                   `out_vars` longtext COLLATE utf8_unicode_ci COMMENT '提供变量',
                                   `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                   `deleted_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
                                   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='变量提供';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `app_api_changelog`
--

DROP TABLE IF EXISTS `app_api_changelog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_api_changelog` (
                                     `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                     `app` varchar(255) DEFAULT NULL COMMENT '所属应用',
                                     `curApiSum` int(11) DEFAULT NULL COMMENT '当前接口总数',
                                     `existApiSum` int(11) DEFAULT NULL COMMENT '已存在接口数',
                                     `newApiSum` int(11) DEFAULT NULL COMMENT '新增接口数',
                                     `deletedApiSum` int(11) DEFAULT NULL COMMENT '删除接口数',
                                     `changedApiSum` int(11) DEFAULT NULL COMMENT '变更接口数',
                                     `checkFailApiSum` int(11) DEFAULT NULL COMMENT '规范检查失败接口数',
                                     `newApiContent` longtext COMMENT '新增接口详情',
                                     `deletedApiContent` longtext COMMENT '删除接口详情',
                                     `changedApiContent` longtext COMMENT '变更接口详情',
                                     `apiCheckResult` varchar(255) DEFAULT NULL COMMENT '接口检查结果',
                                     `apiCheckFailContent` longtext COMMENT '规范检查失败接口详情',
                                     `branch` varchar(255) DEFAULT NULL COMMENT '版本分支',
                                     `remark` varchar(255) DEFAULT NULL COMMENT '备注',
                                     `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                     `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='接口记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `assert_template`
--

DROP TABLE IF EXISTS `assert_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `assert_template` (
                                   `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                   `name` varchar(64) NOT NULL COMMENT '模板名称',
                                   `value` longtext COMMENT '定义信息',
                                   `remark` varchar(255) DEFAULT NULL COMMENT '备注',
                                   `user_name` varchar(100) DEFAULT NULL COMMENT '创建人',
                                   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                   `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`),
                                   UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='断言值多语言模板';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `env_config`
--

DROP TABLE IF EXISTS `env_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `env_config` (
                              `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                              `product` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '产品名称',
                              `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL,
                              `ip` char(39) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '环境域名 / IP / IP:Port',
                              `protocol` enum('http','https') COLLATE utf8_unicode_ci DEFAULT 'http' COMMENT '请求协议',
                              `prepath` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '路由前缀',
                              `threading` enum('yes','no') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'no' COMMENT '是否并发',
                              `thread_number` int(11) DEFAULT NULL COMMENT '并发数',
                              `auth` longtext COLLATE utf8_unicode_ci COMMENT '鉴权信息',
                              `testmode` enum('custom','fuzzing','all') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'fuzzing' COMMENT '测试模式',
                              `swagger_path` varchar(256) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT 'Swagger文档路径',
                              `remark` longtext COLLATE utf8_unicode_ci COMMENT '备注',
                              `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                              `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                              PRIMARY KEY (`id`),
                              UNIQUE KEY `env_config_app_uindex` (`app`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='环境配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `filemanager_setting`
--

DROP TABLE IF EXISTS `filemanager_setting`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `filemanager_setting` (
                                       `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                       `key` varchar(100) DEFAULT NULL,
                                       `value` text,
                                       `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                       `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                       PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `fuzzing_definition`
--

DROP TABLE IF EXISTS `fuzzing_definition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fuzzing_definition` (
                                      `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                      `name` varchar(64) COLLATE utf8_unicode_ci NOT NULL COMMENT '数据描述',
                                      `value` longtext COLLATE utf8_unicode_ci NOT NULL COMMENT '值, e.g.: string, int, bool, list, dict',
                                      `type` enum('string','int','bool','list','dict') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'string' COMMENT '值类型',
                                      `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                      `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                      `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                      PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='随机数据定义';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_menu`
--

DROP TABLE IF EXISTS `goadmin_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_menu` (
                                `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                `parent_id` int(10) unsigned NOT NULL DEFAULT '0',
                                `type` tinyint(3) unsigned NOT NULL DEFAULT '0',
                                `order` int(10) unsigned NOT NULL DEFAULT '0',
                                `title` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                                `icon` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                                `uri` varchar(3000) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
                                `header` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
                                `plugin_name` varchar(150) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
                                `uuid` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
                                `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2578 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_operation_log`
--

DROP TABLE IF EXISTS `goadmin_operation_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_operation_log` (
                                         `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                         `user_id` int(10) unsigned NOT NULL,
                                         `path` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
                                         `method` varchar(10) COLLATE utf8_unicode_ci NOT NULL,
                                         `ip` varchar(15) COLLATE utf8_unicode_ci NOT NULL,
                                         `input` text COLLATE utf8_unicode_ci NOT NULL,
                                         `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                         `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                         PRIMARY KEY (`id`),
                                         KEY `admin_operation_log_user_id_index` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_permissions`
--

DROP TABLE IF EXISTS `goadmin_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_permissions` (
                                       `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                       `name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                                       `slug` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                                       `http_method` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
                                       `http_path` text COLLATE utf8_unicode_ci NOT NULL,
                                       `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                       `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                       PRIMARY KEY (`id`),
                                       UNIQUE KEY `admin_permissions_name_unique` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=236 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_role_menu`
--

DROP TABLE IF EXISTS `goadmin_role_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_role_menu` (
                                     `role_id` int(10) unsigned NOT NULL,
                                     `menu_id` int(10) unsigned NOT NULL,
                                     `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                     `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                     KEY `admin_role_menu_role_id_menu_id_index` (`role_id`,`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_role_permissions`
--

DROP TABLE IF EXISTS `goadmin_role_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_role_permissions` (
                                            `role_id` int(10) unsigned NOT NULL,
                                            `permission_id` int(10) unsigned NOT NULL,
                                            `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                            `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                            UNIQUE KEY `admin_role_permissions` (`role_id`,`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_role_users`
--

DROP TABLE IF EXISTS `goadmin_role_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_role_users` (
                                      `role_id` int(10) unsigned NOT NULL,
                                      `user_id` int(10) unsigned NOT NULL,
                                      `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                      `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                      UNIQUE KEY `admin_user_roles` (`role_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_roles`
--

DROP TABLE IF EXISTS `goadmin_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_roles` (
                                 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                 `name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                                 `slug` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
                                 `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                 `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `admin_roles_name_unique` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_session`
--

DROP TABLE IF EXISTS `goadmin_session`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_session` (
                                   `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                   `sid` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
                                   `values` varchar(3000) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
                                   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                   PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_site`
--

DROP TABLE IF EXISTS `goadmin_site`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_site` (
                                `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                `key` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
                                `value` longtext COLLATE utf8_unicode_ci,
                                `description` varchar(3000) COLLATE utf8_unicode_ci DEFAULT NULL,
                                `state` tinyint(3) unsigned NOT NULL DEFAULT '0',
                                `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=76 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_user_permissions`
--

DROP TABLE IF EXISTS `goadmin_user_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_user_permissions` (
                                            `user_id` int(10) unsigned NOT NULL,
                                            `permission_id` int(10) unsigned NOT NULL,
                                            `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                            `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                            UNIQUE KEY `admin_user_permissions` (`user_id`,`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goadmin_users`
--

DROP TABLE IF EXISTS `goadmin_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goadmin_users` (
                                 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                 `username` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
                                 `password` varchar(100) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
                                 `name` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
                                 `avatar` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
                                 `remember_token` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
                                 `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                 `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `admin_users_username_unique` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `parameter_definition`
--

DROP TABLE IF EXISTS `parameter_definition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `parameter_definition` (
                                        `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                        `name` varchar(64) COLLATE utf8_unicode_ci NOT NULL COMMENT '参数名称/接口ID',
                                        `value` longtext COLLATE utf8_unicode_ci NOT NULL COMMENT '值, e.g.: string, list, dict',
                                        `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '关联产品',
                                        `remark` longtext COLLATE utf8_unicode_ci COMMENT '备注',
                                        `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                        `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                        `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                        PRIMARY KEY (`id`),
                                        UNIQUE KEY `name_app` (`name`,`app`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='参数定义';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `playbook`
--

DROP TABLE IF EXISTS `playbook`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `playbook` (
                            `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                            `name` varchar(64) COLLATE utf8_unicode_ci NOT NULL COMMENT '场景描述',
                            `data_number` text COLLATE utf8_unicode_ci COMMENT '数据序号/标签',
                            `api_list` longtext COLLATE utf8_unicode_ci COMMENT 'API列表',
                            `last_file` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '最近数据文件',
                            `scene_type` enum('1','2') COLLATE utf8_unicode_ci DEFAULT '1' COMMENT '场景类型，1:默认类型，2:比较类型',
                            `priority` int(11) DEFAULT NULL COMMENT '优先级',
                            `run_time` int(11) DEFAULT NULL COMMENT '执行次数',
                            `result` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '测试结果',
                            `fail_reason` longtext COLLATE utf8_unicode_ci COMMENT '失败原因',
                            `remark` longtext COLLATE utf8_unicode_ci COMMENT '备注',
                            `user_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '创建人',
                            `product` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属产品',
                            `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                            `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `name_product` (`name`,`product`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='场景列表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `product`
--

DROP TABLE IF EXISTS `product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product` (
                           `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                           `product` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
                           `ip` char(39) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '环境域名 / IP / IP:Port',
                           `protocol` enum('http','https') COLLATE utf8mb4_unicode_ci DEFAULT 'http' COMMENT '请求协议',
                           `threading` enum('yes','no') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'no' COMMENT '是否并发',
                           `thread_number` int(11) DEFAULT NULL COMMENT '并发数',
                           `auth` longtext COLLATE utf8mb4_unicode_ci COMMENT '鉴权信息',
                           `testmode` enum('custom','fuzzing','all') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'fuzzing' COMMENT '测试模式',
                           `apps` longtext COLLATE utf8mb4_unicode_ci COMMENT '关联应用',
                           `env_type` int(11) DEFAULT NULL COMMENT '环境类型，e.g.: 1: 开发，2: 测试，3: 预发，4: 演示，5：生产 ',
                           `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `private_parameter` longtext COLLATE utf8mb4_unicode_ci COMMENT '专用参数',
                           `remark` longtext COLLATE utf8mb4_unicode_ci COMMENT '备注',
                           `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                           `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `product` (`product`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `product_count`
--

DROP TABLE IF EXISTS `product_count`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_count` (
                                 `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                 `all_count` int(11) DEFAULT NULL COMMENT '接口总数',
                                 `automatable_count` int(11) DEFAULT NULL COMMENT '可自动化数',
                                 `unautomatable_count` int(11) DEFAULT NULL COMMENT '不可自动化数',
                                 `auto_test_count` int(11) DEFAULT NULL COMMENT '自动化测试总数',
                                 `untest_count` int(11) DEFAULT NULL COMMENT '未测试总数',
                                 `pass_count` int(11) DEFAULT NULL COMMENT '通过总数',
                                 `fail_count` int(11) DEFAULT NULL COMMENT '失败总数',
                                 `auto_per` double DEFAULT NULL COMMENT '自动化率',
                                 `pass_per` double DEFAULT NULL COMMENT '通过率',
                                 `fail_per` double DEFAULT NULL COMMENT '失败率',
                                 `product` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL,
                                 `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                 `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='产品统计';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `scene_data`
--

DROP TABLE IF EXISTS `scene_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scene_data` (
                              `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                              `name` varchar(64) COLLATE utf8_unicode_ci NOT NULL COMMENT '数据描述',
                              `api_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '接口ID',
                              `app` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '所属应用',
                              `file_name` longtext COLLATE utf8_unicode_ci COMMENT '文件名称',
                              `content` longtext COLLATE utf8_unicode_ci COMMENT '文件内容',
                              `run_time` int(11) DEFAULT NULL COMMENT '执行次数',
                              `result` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '测试结果',
                              `fail_reason` longtext COLLATE utf8_unicode_ci COMMENT '失败原因',
                              `remark` longtext COLLATE utf8_unicode_ci COMMENT '备注',
                              `user_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '创建人',
                              `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                              `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                              PRIMARY KEY (`id`),
                              UNIQUE KEY `scene_data_pk_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='场景数据';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `scene_data_test_history`
--

DROP TABLE IF EXISTS `scene_data_test_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scene_data_test_history` (
                                           `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                           `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '数据描述',
                                           `api_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '接口ID',
                                           `app` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所属应用',
                                           `content` longtext COLLATE utf8mb4_unicode_ci COMMENT '.yml文件内容',
                                           `result` varchar(6) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '测试结果',
                                           `fail_reason` longtext COLLATE utf8mb4_unicode_ci COMMENT '失败原因',
                                           `env_type` int(11) DEFAULT NULL COMMENT '环境类型，e.g.: 1: 开发，2: 测试，3: 预发，4: 演示，5：生产',
                                           `product` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                           `remark` longtext COLLATE utf8mb4_unicode_ci COMMENT '备注',
                                           `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                           `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                           `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                           PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='场景数据测试记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `scene_test_history`
--

DROP TABLE IF EXISTS `scene_test_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scene_test_history` (
                                      `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                      `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '场景描述',
                                      `api_list` longtext COLLATE utf8mb4_unicode_ci COMMENT 'API列表',
                                      `last_file` longtext COLLATE utf8mb4_unicode_ci COMMENT '最近数据文件',
                                      `scene_type` enum('1','2') COLLATE utf8mb4_unicode_ci DEFAULT '1' COMMENT '场景类型，1:默认类型，2:比较类型',
                                      `result` varchar(6) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '测试结果',
                                      `fail_reason` longtext COLLATE utf8mb4_unicode_ci COMMENT '失败原因',
                                      `product` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所属产品',
                                      `env_type` int(11) DEFAULT NULL COMMENT '环境类型，e.g.: 1: 开发，2: 测试，3: 预发，4: 演示，5：生产',
                                      `remark` longtext COLLATE utf8mb4_unicode_ci COMMENT '备注',
                                      `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                      `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                      `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                      PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='场景测试记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `schedule`
--

DROP TABLE IF EXISTS `schedule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `schedule` (
                            `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                            `task_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务名称',
                            `task_mode` enum('cron','once','day','week') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'once' COMMENT '任务模式',
                            `crontab` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'cron表达式',
                            `threading` enum('yes','no') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'no' COMMENT '是否并发',
                            `task_type` enum('data','scene') COLLATE utf8mb4_unicode_ci DEFAULT 'data' COMMENT '任务类型',
                            `data_number` text COLLATE utf8mb4_unicode_ci COMMENT '数据序号/标签',
                            `data_list` longtext COLLATE utf8mb4_unicode_ci COMMENT '关联数据',
                            `scene_number` text COLLATE utf8mb4_unicode_ci COMMENT '场景序号/标签',
                            `scene_list` longtext COLLATE utf8mb4_unicode_ci COMMENT '关联场景',
                            `product_list` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联产品',
                            `task_status` enum('running','stopped','finished','not_started') COLLATE utf8mb4_unicode_ci DEFAULT 'not_started' COMMENT '任务状态',
                            `time4week` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '每时',
                            `time4day` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '每时',
                            `week` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '每周',
                            `remark` longtext COLLATE utf8mb4_unicode_ci COMMENT '备注',
                            `last_at` timestamp NULL DEFAULT NULL COMMENT '上次执行时间',
                            `next_at` timestamp NULL DEFAULT NULL COMMENT '下次执行时间',
                            `user_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '创建人',
                            `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                            `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `task_name` (`task_name`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='定时任务';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sys_parameter`
--

DROP TABLE IF EXISTS `sys_parameter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sys_parameter` (
                                 `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                 `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '参数名称',
                                 `value_list` longtext COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '枚举列表',
                                 `remark` longtext COLLATE utf8mb4_unicode_ci COMMENT '备注',
                                 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                                 `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统参数';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `test_case`
--

DROP TABLE IF EXISTS `test_case`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `test_case` (
                             `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                             `case_number` varchar(50) NOT NULL COMMENT '用例编号',
                             `case_name` varchar(255) DEFAULT NULL COMMENT '用例名称',
                             `case_type` varchar(16) DEFAULT NULL COMMENT '用例类型',
                             `priority` varchar(6) DEFAULT NULL COMMENT '优先级',
                             `pre_condition` longtext COMMENT '前置条件',
                             `test_range` longtext COMMENT '测试范围',
                             `test_steps` longtext COMMENT '测试步骤',
                             `expect_result` longtext COMMENT '预期结果',
                             `auto` varchar(3) DEFAULT NULL COMMENT '是否自动化',
                             `fun_developer` varchar(10) DEFAULT NULL COMMENT '功能开发者',
                             `case_designer` varchar(10) DEFAULT NULL COMMENT '用例设计者',
                             `case_executor` varchar(10) DEFAULT NULL COMMENT '用例执行者',
                             `test_time` varchar(12) DEFAULT NULL COMMENT '测试时间',
                             `test_result` varchar(12) DEFAULT NULL COMMENT '测试结果',
                             `module` varchar(255) DEFAULT NULL COMMENT '所属模块',
                             `remark` varchar(255) DEFAULT NULL COMMENT '备注',
                             `scene` varchar(255) DEFAULT NULL COMMENT '关联场景',
                             `product` varchar(255) DEFAULT NULL COMMENT '所属产品',
                             `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                             `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `case_number_product` (`case_number`,`product`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='测试用例';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-01-05 20:59:06

-- MySQL dump 10.14  Distrib 5.5.60-MariaDB, for Linux (x86_64)
--
-- Host: 10.57.242.40    Database: data4test
-- ------------------------------------------------------
-- Server version	5.7.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `filemanager_setting`
--

LOCK TABLES `filemanager_setting` WRITE;
/*!40000 ALTER TABLE `filemanager_setting` DISABLE KEYS */;
INSERT INTO `filemanager_setting` VALUES (1,'roots','{\"api\":{\"Path\":\"mgmt/api\",\"Title\":\"接口文件\"},\"case\":{\"Path\":\"mgmt/case\",\"Title\":\"用例文件\"},\"common\":{\"Path\":\"mgmt/common\",\"Title\":\"公用文件(勿删)\"},\"data\":{\"Path\":\"mgmt/data\",\"Title\":\"数据文件\"},\"download\":{\"Path\":\"mgmt/download\",\"Title\":\"下载文件\"},\"history\":{\"Path\":\"mgmt/history\",\"Title\":\"历史记录\"},\"log\":{\"Path\":\"mgmt/log\",\"Title\":\"日志管理\"},\"old\":{\"Path\":\"mgmt/old\",\"Title\":\"历史版本\"},\"upload\":{\"Path\":\"mgmt/upload\",\"Title\":\"上传文件\"}}','2020-12-07 02:35:59','2020-12-07 02:35:59'),(2,'allowMove','1','2020-12-07 02:36:00','2020-12-07 02:36:00'),(3,'conn','default','2020-12-07 02:36:00','2020-12-07 02:36:00'),(4,'allowUpload','1','2020-12-07 02:36:00','2020-12-07 02:36:00'),(5,'allowDelete','1','2020-12-07 02:36:00','2020-12-07 02:36:00'),(6,'allowRename','1','2020-12-07 02:36:00','2020-12-07 02:36:00'),(7,'allowDownload','1','2020-12-07 02:36:00','2020-12-07 02:36:00'),(8,'allowCreateDir','1','2020-12-07 02:36:00','2020-12-07 02:36:00');
/*!40000 ALTER TABLE `filemanager_setting` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goadmin_menu`
--

LOCK TABLES `goadmin_menu` WRITE;
/*!40000 ALTER TABLE `goadmin_menu` DISABLE KEYS */;
INSERT INTO `goadmin_menu` VALUES (1,0,1,40,'Admin','fa-tasks','',NULL,'',NULL,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(2,1,1,40,'Users','fa-users','/info/manager',NULL,'',NULL,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(3,1,1,41,'Roles','fa-user','/info/roles',NULL,'',NULL,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(4,1,1,42,'Permission','fa-ban','/info/permission',NULL,'',NULL,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(5,1,1,43,'Menu','fa-bars','/menu',NULL,'',NULL,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(6,50,1,39,'Operation log','fa-history','/info/op','','',NULL,'2019-09-09 16:00:00','2021-03-18 14:52:30'),(7,0,1,1,'Dashboard','fa-bar-chart','/dashboard','','',NULL,'2019-09-09 16:00:00','2020-12-28 14:16:17'),(8,0,0,16,'环境','fa-cog','','','',NULL,'2020-11-23 02:47:56','2020-11-23 02:57:07'),(9,8,0,17,'应用配置','fa-bars','/info/env_config','','',NULL,'2020-11-23 02:57:40','2021-07-29 06:04:57'),(11,8,0,19,'参数定义','fa-bars','/info/parameter_definition','','',NULL,'2020-11-23 02:58:49','2021-03-22 09:49:42'),(12,0,0,3,'接口','fa-bank','','','',NULL,'2020-11-23 02:59:22','2021-03-21 07:08:41'),(14,0,0,11,'结果','fa-cubes','','','',NULL,'2020-11-23 02:59:57','2020-11-23 02:59:57'),(17,14,0,14,'变量提供','fa-bars','/info/api_test_result','','',NULL,'2020-11-23 03:01:36','2021-03-30 03:02:07'),(18,14,0,13,'结果详情','fa-bars','/info/api_test_detail','','',NULL,'2020-11-23 03:02:09','2020-11-23 03:06:44'),(22,12,0,5,'接口关系','fa-bars','/info/api_relation','','',NULL,'2020-11-23 03:03:48','2021-03-22 09:47:56'),(25,12,0,6,'接口定义','fa-bars','/info/api_definition','','',NULL,'2020-11-24 09:22:39','2021-03-22 09:48:08'),(34,28,0,22,'Swagger文件','fa-bars','/fm/api/list','','',NULL,'2020-12-07 03:06:28','2020-12-07 03:06:28'),(37,28,0,21,'用例文件','fa-bars','/fm/test/list','','',NULL,'2020-12-07 03:08:41','2020-12-07 03:08:41'),(38,28,0,20,'公共文件','fa-bars','/fm/file/list','','',NULL,'2020-12-07 03:09:05','2020-12-07 03:09:05'),(39,0,0,22,'文件','fa-files-o','','','',NULL,'2020-12-07 03:20:35','2020-12-10 03:17:02'),(40,39,0,29,'公共文件','fa-bars','/fm/common/list','','',NULL,'2020-12-07 03:23:33','2023-08-11 10:05:23'),(41,39,0,28,'用例文件','fa-bars','/fm/case/list','','',NULL,'2020-12-07 03:23:49','2023-08-11 10:05:06'),(42,39,0,27,'API文件','fa-bars','/fm/api/list','','',NULL,'2020-12-07 03:24:04','2020-12-07 03:24:04'),(43,12,0,3,'测试数据','fa-bars','/info/api_test_data','','',NULL,'2020-12-08 08:20:47','2020-12-08 08:20:47'),(44,39,0,30,'日志文件','fa-bars','/fm/log/list','','',NULL,'2020-12-10 03:08:30','2020-12-10 03:17:15'),(45,0,0,9,'日志文件','fa-bars','/fm/log/list','','filemanager',NULL,'2020-12-10 03:09:23','2020-12-10 03:16:07'),(46,0,0,7,'用例文件','fa-bars','/fm/case/list','','filemanager',NULL,'2020-12-10 03:10:06','2023-08-11 10:00:56'),(47,0,0,6,'API文件','fa-bars','/fm/api/list','','filemanager',NULL,'2020-12-10 03:10:35','2020-12-10 03:10:35'),(48,0,0,8,'公用文件','fa-bars','/fm/common/list','','filemanager',NULL,'2020-12-10 03:10:56','2023-08-11 10:01:05'),(50,0,0,39,'日志','fa-500px','','','',NULL,'2021-03-18 14:52:09','2021-03-18 14:53:17'),(51,12,0,4,'模糊数据','fa-bars','/info/api_fuzzing_data','','',NULL,'2021-03-21 07:07:45','2021-03-22 11:26:03'),(101,0,0,31,'说明','fa-bars','','','',NULL,'2021-03-23 06:14:06','2021-03-23 06:16:25'),(102,101,0,33,'系统蓝图','fa-bars','/librarian/plan/blue_print','','',NULL,'2021-03-23 06:15:53','2024-01-05 11:46:06'),(103,101,0,37,'Mock应用说明','fa-bars','/librarian/design/mock_design','','',NULL,'2021-03-23 06:17:34','2024-01-05 11:47:48'),(109,101,0,32,'架构图','fa-bars','/librarian/arch/arch','','',NULL,'2021-03-23 06:29:45','2021-03-24 05:57:14'),(126,101,0,31,'概览','fa-bars','/librarian/README','','',NULL,'2021-03-24 03:27:58','2021-03-24 03:27:58'),(127,101,0,1,'概览','fa-bars','/librarian/README','','librarian',NULL,'2021-03-24 03:30:28','2021-03-24 03:31:42'),(136,8,0,21,'模糊因子','fa-bars','/info/fuzzing_definition','','',NULL,'2021-03-26 01:57:30','2021-03-26 01:57:30'),(137,0,0,9,'场景','fa-bars','','','',NULL,'2021-06-03 03:42:23','2021-06-03 03:42:23'),(138,137,0,9,'场景列表','fa-bars','/info/playbook','','',NULL,'2021-06-03 03:43:06','2021-06-03 03:43:57'),(139,137,0,10,'数据列表','fa-bars','/info/scene_data','','',NULL,'2021-06-03 03:43:19','2021-07-29 06:20:10'),(140,0,0,1,'数据文件','fa-bars','/fm/data/list','','filemanager',NULL,'2021-06-03 06:02:31','2023-08-11 10:00:46'),(141,39,0,22,'数据文件','fa-bars','/fm/data/list','','',NULL,'2021-06-04 02:08:36','2023-08-11 10:04:51'),(202,8,0,16,'产品配置','fa-bars','/info/product','','',NULL,'2021-07-19 07:44:29','2021-07-19 07:44:29'),(203,0,0,15,'任务','fa-arrows-alt','','','',NULL,'2021-07-23 09:09:57','2021-07-29 05:55:14'),(204,203,0,15,'定时任务','fa-bars','/info/schedule','','',NULL,'2021-07-23 09:10:20','2021-07-23 09:10:20'),(205,0,0,2,'上传文件','fa-bars','/fm/upload/list','','filemanager',NULL,'2021-08-04 02:47:48','2021-08-04 02:47:48'),(206,0,0,3,'下载文件','fa-bars','/fm/download/list','','filemanager',NULL,'2021-08-04 02:48:01','2021-08-04 02:48:01'),(207,0,0,4,'历史记录','fa-bars','/fm/history/list','','filemanager',NULL,'2021-08-04 02:48:10','2021-08-04 02:58:05'),(208,39,0,23,'上传文件','fa-bars','/fm/upload/list','','',NULL,'2021-08-04 02:50:17','2021-08-04 02:50:17'),(209,39,0,24,'历史记录','fa-bars','/fm/history/list','','',NULL,'2021-08-04 02:50:27','2021-08-04 02:50:27'),(210,39,0,26,'下载文件','fa-bars','/fm/download/list','','',NULL,'2021-08-04 02:50:35','2021-08-04 02:50:35'),(211,8,0,18,'系统参数','fa-bars','/info/sys_parameter','','',NULL,'2021-08-06 03:55:26','2021-08-06 03:55:26'),(212,14,0,11,'数据测试历史','fa-bars','/info/scene_data_test_history','','',NULL,'2021-08-10 07:03:22','2023-04-27 09:13:17'),(213,14,0,12,'场景测试历史','fa-bars','/info/scene_test_history','','',NULL,'2021-08-10 07:04:29','2022-03-09 05:59:35'),(214,39,0,25,'历史版本','fa-bars','/fm/old/list','','',NULL,'2021-08-30 06:46:40','2021-08-30 06:46:40'),(215,0,0,5,'历史版本','fa-bars','/fm/old/list','','filemanager',NULL,'2021-08-30 06:47:36','2021-08-30 06:47:36'),(216,0,0,8,'用例','fa-bars','','','',NULL,'2021-09-05 02:48:39','2021-09-05 02:48:39'),(217,216,0,8,'测试用例','fa-bars','/info/test_case','','',NULL,'2021-09-05 02:49:29','2021-09-05 02:49:29'),(219,101,0,38,'更新记录','fa-bars','/librarian/update/change_log','','',NULL,'2021-11-08 09:55:24','2021-11-08 09:55:24'),(236,0,0,2,'控制台','fa-bars','/likePostman','','',NULL,'2021-12-03 06:32:52','2022-03-07 02:24:16'),(237,12,0,7,'接口记录表','fa-bars','/info/app_api_changelog','','',NULL,'2022-03-08 02:25:03','2022-03-08 02:25:03'),(1336,101,0,34,'参数应用说明','fa-bars','/librarian/design/parameter_design','','',NULL,'2023-02-23 09:54:47','2024-01-05 11:48:23'),(1367,101,0,35,'断言应用说明','fa-bars','/librarian/design/assert_design','','',NULL,'2023-02-27 11:47:11','2024-01-05 11:46:47'),(1761,101,0,36,'动作应用说明','fa-bars','/librarian/design/action_design','','',NULL,'2023-06-15 08:57:30','2024-01-05 11:47:07'),(2224,8,0,20,'断言值模板','fa-bars','/info/assert_template','','',NULL,'2023-12-11 07:31:19','2023-12-11 08:09:41'),(2564,0,0,1,'概览','fa-file-o','/librarian/README','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2565,0,0,2,'架构图','fa-file-o','/librarian/arch/arch','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2566,0,0,3,'系统蓝图','fa-file-o','/librarian/plan/blue_print','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2567,0,0,4,'场景测试说明','fa-file-o','/librarian/design/scene_design','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2568,0,0,5,'参数应用说明','fa-file-o','/librarian/design/parameter_design','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2569,0,0,6,'动作应用说明','fa-file-o','/librarian/design/action_design','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2570,0,0,7,'断言应用说明','fa-file-o','/librarian/design/assert_design','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2571,0,0,8,'Mock应用说明','fa-file-o','/librarian/design/mock_design','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2572,0,0,9,'模糊测试说明','fa-file-o','/librarian/design/relation_design','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2573,0,0,10,'测试用例说明','fa-file-o','/librarian/design/testcase_design','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2574,0,0,11,'性能测试设计','fa-file-o','/librarian/design/perf_design','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2575,0,0,12,'开发待办','fa-file-o','/librarian/plan/todo','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2576,0,0,13,'开发须知','fa-file-o','/librarian/development/must_know','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11'),(2577,0,0,14,'更新记录','fa-file-o','/librarian/update/change_log','','librarian','','2024-01-05 12:45:11','2024-01-05 12:45:11');
/*!40000 ALTER TABLE `goadmin_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goadmin_permissions`
--

LOCK TABLES `goadmin_permissions` WRITE;
/*!40000 ALTER TABLE `goadmin_permissions` DISABLE KEYS */;
INSERT INTO `goadmin_permissions` VALUES (1,'All permission','*','','*','2019-09-09 16:00:00','2019-09-09 16:00:00'),(2,'Dashboard','dashboard','GET,PUT,POST,DELETE','/','2019-09-09 16:00:00','2019-09-09 16:00:00'),(3,'api_case 查询','api_case_query','GET','/info/api_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(4,'api_case 编辑页显示','api_case_show_edit','GET','/info/api_case/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(5,'api_case 新建记录页显示','api_case_show_create','GET','/info/api_case/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(6,'api_case 编辑','api_case_edit','POST','/edit/api_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(7,'api_case 新建','api_case_create','POST','/new/api_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(8,'api_case 删除','api_case_delete','POST','/delete/api_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(9,'api_case 导出','api_case_export','POST','/export/api_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(10,'api_sum_up 查询','api_sum_up_query','GET','/info/api_sum_up','2020-11-23 02:06:28','2020-11-23 02:06:28'),(11,'api_sum_up 编辑页显示','api_sum_up_show_edit','GET','/info/api_sum_up/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(12,'api_sum_up 新建记录页显示','api_sum_up_show_create','GET','/info/api_sum_up/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(13,'api_sum_up 编辑','api_sum_up_edit','POST','/edit/api_sum_up','2020-11-23 02:06:28','2020-11-23 02:06:28'),(14,'api_sum_up 新建','api_sum_up_create','POST','/new/api_sum_up','2020-11-23 02:06:28','2020-11-23 02:06:28'),(15,'api_sum_up 删除','api_sum_up_delete','POST','/delete/api_sum_up','2020-11-23 02:06:28','2020-11-23 02:06:28'),(16,'api_sum_up 导出','api_sum_up_export','POST','/export/api_sum_up','2020-11-23 02:06:28','2020-11-23 02:06:28'),(17,'api_test_detail 查询','api_test_detail_query','GET','/info/api_test_detail','2020-11-23 02:06:28','2020-11-23 02:06:28'),(18,'api_test_detail 编辑页显示','api_test_detail_show_edit','GET','/info/api_test_detail/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(19,'api_test_detail 新建记录页显示','api_test_detail_show_create','GET','/info/api_test_detail/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(20,'api_test_detail 编辑','api_test_detail_edit','POST','/edit/api_test_detail','2020-11-23 02:06:28','2020-11-23 02:06:28'),(21,'api_test_detail 新建','api_test_detail_create','POST','/new/api_test_detail','2020-11-23 02:06:28','2020-11-23 02:06:28'),(22,'api_test_detail 删除','api_test_detail_delete','POST','/delete/api_test_detail','2020-11-23 02:06:28','2020-11-23 02:06:28'),(23,'api_test_detail 导出','api_test_detail_export','POST','/export/api_test_detail','2020-11-23 02:06:28','2020-11-23 02:06:28'),(24,'api_test_result 查询','api_test_result_query','GET','/info/api_test_result','2020-11-23 02:06:28','2020-11-23 02:06:28'),(25,'api_test_result 编辑页显示','api_test_result_show_edit','GET','/info/api_test_result/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(26,'api_test_result 新建记录页显示','api_test_result_show_create','GET','/info/api_test_result/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(27,'api_test_result 编辑','api_test_result_edit','POST','/edit/api_test_result','2020-11-23 02:06:28','2020-11-23 02:06:28'),(28,'api_test_result 新建','api_test_result_create','POST','/new/api_test_result','2020-11-23 02:06:28','2020-11-23 02:06:28'),(29,'api_test_result 删除','api_test_result_delete','POST','/delete/api_test_result','2020-11-23 02:06:28','2020-11-23 02:06:28'),(30,'api_test_result 导出','api_test_result_export','POST','/export/api_test_result','2020-11-23 02:06:28','2020-11-23 02:06:28'),(31,'case_test_count 查询','case_test_count_query','GET','/info/case_test_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(32,'case_test_count 编辑页显示','case_test_count_show_edit','GET','/info/case_test_count/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(33,'case_test_count 新建记录页显示','case_test_count_show_create','GET','/info/case_test_count/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(34,'case_test_count 编辑','case_test_count_edit','POST','/edit/case_test_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(35,'case_test_count 新建','case_test_count_create','POST','/new/case_test_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(36,'case_test_count 删除','case_test_count_delete','POST','/delete/case_test_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(37,'case_test_count 导出','case_test_count_export','POST','/export/case_test_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(38,'common_variable 查询','common_variable_query','GET','/info/common_variable','2020-11-23 02:06:28','2020-11-23 02:06:28'),(39,'common_variable 编辑页显示','common_variable_show_edit','GET','/info/common_variable/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(40,'common_variable 新建记录页显示','common_variable_show_create','GET','/info/common_variable/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(41,'common_variable 编辑','common_variable_edit','POST','/edit/common_variable','2020-11-23 02:06:28','2020-11-23 02:06:28'),(42,'common_variable 新建','common_variable_create','POST','/new/common_variable','2020-11-23 02:06:28','2020-11-23 02:06:28'),(43,'common_variable 删除','common_variable_delete','POST','/delete/common_variable','2020-11-23 02:06:28','2020-11-23 02:06:28'),(44,'common_variable 导出','common_variable_export','POST','/export/common_variable','2020-11-23 02:06:28','2020-11-23 02:06:28'),(45,'host 查询','host_query','GET','/info/host','2020-11-23 02:06:28','2020-11-23 02:06:28'),(46,'host 编辑页显示','host_show_edit','GET','/info/host/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(47,'host 新建记录页显示','host_show_create','GET','/info/host/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(48,'host 编辑','host_edit','POST','/edit/host','2020-11-23 02:06:28','2020-11-23 02:06:28'),(49,'host 新建','host_create','POST','/new/host','2020-11-23 02:06:28','2020-11-23 02:06:28'),(50,'host 删除','host_delete','POST','/delete/host','2020-11-23 02:06:28','2020-11-23 02:06:28'),(51,'host 导出','host_export','POST','/export/host','2020-11-23 02:06:28','2020-11-23 02:06:28'),(52,'issue 查询','issue_query','GET','/info/issue','2020-11-23 02:06:28','2020-11-23 02:06:28'),(53,'issue 编辑页显示','issue_show_edit','GET','/info/issue/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(54,'issue 新建记录页显示','issue_show_create','GET','/info/issue/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(55,'issue 编辑','issue_edit','POST','/edit/issue','2020-11-23 02:06:28','2020-11-23 02:06:28'),(56,'issue 新建','issue_create','POST','/new/issue','2020-11-23 02:06:28','2020-11-23 02:06:28'),(57,'issue 删除','issue_delete','POST','/delete/issue','2020-11-23 02:06:28','2020-11-23 02:06:28'),(58,'issue 导出','issue_export','POST','/export/issue','2020-11-23 02:06:28','2020-11-23 02:06:28'),(59,'issue_milestone_count 查询','issue_milestone_count_query','GET','/info/issue_milestone_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(60,'issue_milestone_count 编辑页显示','issue_milestone_count_show_edit','GET','/info/issue_milestone_count/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(61,'issue_milestone_count 新建记录页显示','issue_milestone_count_show_create','GET','/info/issue_milestone_count/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(62,'issue_milestone_count 编辑','issue_milestone_count_edit','POST','/edit/issue_milestone_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(63,'issue_milestone_count 新建','issue_milestone_count_create','POST','/new/issue_milestone_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(64,'issue_milestone_count 删除','issue_milestone_count_delete','POST','/delete/issue_milestone_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(65,'issue_milestone_count 导出','issue_milestone_count_export','POST','/export/issue_milestone_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(66,'issue_tag_count 查询','issue_tag_count_query','GET','/info/issue_tag_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(67,'issue_tag_count 编辑页显示','issue_tag_count_show_edit','GET','/info/issue_tag_count/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(68,'issue_tag_count 新建记录页显示','issue_tag_count_show_create','GET','/info/issue_tag_count/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(69,'issue_tag_count 编辑','issue_tag_count_edit','POST','/edit/issue_tag_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(70,'issue_tag_count 新建','issue_tag_count_create','POST','/new/issue_tag_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(71,'issue_tag_count 删除','issue_tag_count_delete','POST','/delete/issue_tag_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(72,'issue_tag_count 导出','issue_tag_count_export','POST','/export/issue_tag_count','2020-11-23 02:06:28','2020-11-23 02:06:28'),(73,'product_gitlab 查询','product_gitlab_query','GET','/info/product_gitlab','2020-11-23 02:06:28','2020-11-23 02:06:28'),(74,'product_gitlab 编辑页显示','product_gitlab_show_edit','GET','/info/product_gitlab/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(75,'product_gitlab 新建记录页显示','product_gitlab_show_create','GET','/info/product_gitlab/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(76,'product_gitlab 编辑','product_gitlab_edit','POST','/edit/product_gitlab','2020-11-23 02:06:28','2020-11-23 02:06:28'),(77,'product_gitlab 新建','product_gitlab_create','POST','/new/product_gitlab','2020-11-23 02:06:28','2020-11-23 02:06:28'),(78,'product_gitlab 删除','product_gitlab_delete','POST','/delete/product_gitlab','2020-11-23 02:06:28','2020-11-23 02:06:28'),(79,'product_gitlab 导出','product_gitlab_export','POST','/export/product_gitlab','2020-11-23 02:06:28','2020-11-23 02:06:28'),(80,'test_case 查询','test_case_query','GET','/info/test_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(81,'test_case 编辑页显示','test_case_show_edit','GET','/info/test_case/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(82,'test_case 新建记录页显示','test_case_show_create','GET','/info/test_case/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(83,'test_case 编辑','test_case_edit','POST','/edit/test_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(84,'test_case 新建','test_case_create','POST','/new/test_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(85,'test_case 删除','test_case_delete','POST','/delete/test_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(86,'test_case 导出','test_case_export','POST','/export/test_case','2020-11-23 02:06:28','2020-11-23 02:06:28'),(87,'test_progress_schedule 查询','test_progress_schedule_query','GET','/info/test_progress_schedule','2020-11-23 02:06:28','2020-11-23 02:06:28'),(88,'test_progress_schedule 编辑页显示','test_progress_schedule_show_edit','GET','/info/test_progress_schedule/edit','2020-11-23 02:06:28','2020-11-23 02:06:28'),(89,'test_progress_schedule 新建记录页显示','test_progress_schedule_show_create','GET','/info/test_progress_schedule/new','2020-11-23 02:06:28','2020-11-23 02:06:28'),(90,'test_progress_schedule 编辑','test_progress_schedule_edit','POST','/edit/test_progress_schedule','2020-11-23 02:06:28','2020-11-23 02:06:28'),(91,'test_progress_schedule 新建','test_progress_schedule_create','POST','/new/test_progress_schedule','2020-11-23 02:06:28','2020-11-23 02:06:28'),(92,'test_progress_schedule 删除','test_progress_schedule_delete','POST','/delete/test_progress_schedule','2020-11-23 02:06:28','2020-11-23 02:06:28'),(93,'test_progress_schedule 导出','test_progress_schedule_export','POST','/export/test_progress_schedule','2020-11-23 02:06:28','2020-11-23 02:06:28'),(94,'api_detail 查询','api_detail_query','GET','/info/api_detail','2020-11-24 09:18:47','2020-11-24 09:18:47'),(95,'api_detail 编辑页显示','api_detail_show_edit','GET','/info/api_detail/edit','2020-11-24 09:18:47','2020-11-24 09:18:47'),(96,'api_detail 新建记录页显示','api_detail_show_create','GET','/info/api_detail/new','2020-11-24 09:18:47','2020-11-24 09:18:47'),(97,'api_detail 编辑','api_detail_edit','POST','/edit/api_detail','2020-11-24 09:18:47','2020-11-24 09:18:47'),(98,'api_detail 新建','api_detail_create','POST','/new/api_detail','2020-11-24 09:18:47','2020-11-24 09:18:47'),(99,'api_detail 删除','api_detail_delete','POST','/delete/api_detail','2020-11-24 09:18:47','2020-11-24 09:18:47'),(100,'api_detail 导出','api_detail_export','POST','/export/api_detail','2020-11-24 09:18:47','2020-11-24 09:18:47'),(101,'api_test_data 查询','api_test_data_query','GET','/info/api_test_data','2020-12-08 08:18:54','2020-12-08 08:18:54'),(102,'api_test_data 编辑页显示','api_test_data_show_edit','GET','/info/api_test_data/edit','2020-12-08 08:18:54','2020-12-08 08:18:54'),(103,'api_test_data 新建记录页显示','api_test_data_show_create','GET','/info/api_test_data/new','2020-12-08 08:18:54','2020-12-08 08:18:54'),(104,'api_test_data 编辑','api_test_data_edit','POST','/edit/api_test_data','2020-12-08 08:18:54','2020-12-08 08:18:54'),(105,'api_test_data 新建','api_test_data_create','POST','/new/api_test_data','2020-12-08 08:18:54','2020-12-08 08:18:54'),(106,'api_test_data 删除','api_test_data_delete','POST','/delete/api_test_data','2020-12-08 08:18:54','2020-12-08 08:18:54'),(107,'api_test_data 导出','api_test_data_export','POST','/export/api_test_data','2020-12-08 08:18:54','2020-12-08 08:18:54'),(108,'testcase_count 查询','testcase_count_query','GET','/info/testcase_count','2020-12-15 03:17:56','2020-12-15 03:17:56'),(109,'testcase_count 编辑页显示','testcase_count_show_edit','GET','/info/testcase_count/edit','2020-12-15 03:17:56','2020-12-15 03:17:56'),(110,'testcase_count 新建记录页显示','testcase_count_show_create','GET','/info/testcase_count/new','2020-12-15 03:17:56','2020-12-15 03:17:56'),(111,'testcase_count 编辑','testcase_count_edit','POST','/edit/testcase_count','2020-12-15 03:17:56','2020-12-15 03:17:56'),(112,'testcase_count 新建','testcase_count_create','POST','/new/testcase_count','2020-12-15 03:17:56','2020-12-15 03:17:56'),(113,'testcase_count 删除','testcase_count_delete','POST','/delete/testcase_count','2020-12-15 03:17:56','2020-12-15 03:17:56'),(114,'testcase_count 导出','testcase_count_export','POST','/export/testcase_count','2020-12-15 03:17:56','2020-12-15 03:17:56'),(115,'env_config 查询','env_config_query','GET','/info/env_config','2021-03-22 09:16:04','2021-03-22 09:16:04'),(116,'env_config 编辑页显示','env_config_show_edit','GET','/info/env_config/edit','2021-03-22 09:16:04','2021-03-22 09:16:04'),(117,'env_config 新建记录页显示','env_config_show_create','GET','/info/env_config/new','2021-03-22 09:16:04','2021-03-22 09:16:04'),(118,'env_config 编辑','env_config_edit','POST','/edit/env_config','2021-03-22 09:16:04','2021-03-22 09:16:04'),(119,'env_config 新建','env_config_create','POST','/new/env_config','2021-03-22 09:16:04','2021-03-22 09:16:04'),(120,'env_config 删除','env_config_delete','POST','/delete/env_config','2021-03-22 09:16:04','2021-03-22 09:16:04'),(121,'env_config 导出','env_config_export','POST','/export/env_config','2021-03-22 09:16:04','2021-03-22 09:16:04'),(122,'api_definition 查询','api_definition_query','GET','/info/api_definition','2021-03-22 09:23:54','2021-03-22 09:23:54'),(123,'api_definition 编辑页显示','api_definition_show_edit','GET','/info/api_definition/edit','2021-03-22 09:23:54','2021-03-22 09:23:54'),(124,'api_definition 新建记录页显示','api_definition_show_create','GET','/info/api_definition/new','2021-03-22 09:23:54','2021-03-22 09:23:54'),(125,'api_definition 编辑','api_definition_edit','POST','/edit/api_definition','2021-03-22 09:23:54','2021-03-22 09:23:54'),(126,'api_definition 新建','api_definition_create','POST','/new/api_definition','2021-03-22 09:23:54','2021-03-22 09:23:54'),(127,'api_definition 删除','api_definition_delete','POST','/delete/api_definition','2021-03-22 09:23:54','2021-03-22 09:23:54'),(128,'api_definition 导出','api_definition_export','POST','/export/api_definition','2021-03-22 09:23:54','2021-03-22 09:23:54'),(129,'api_relation 查询','api_relation_query','GET','/info/api_relation','2021-03-22 09:29:01','2021-03-22 09:29:01'),(130,'api_relation 编辑页显示','api_relation_show_edit','GET','/info/api_relation/edit','2021-03-22 09:29:01','2021-03-22 09:29:01'),(131,'api_relation 新建记录页显示','api_relation_show_create','GET','/info/api_relation/new','2021-03-22 09:29:01','2021-03-22 09:29:01'),(132,'api_relation 编辑','api_relation_edit','POST','/edit/api_relation','2021-03-22 09:29:01','2021-03-22 09:29:01'),(133,'api_relation 新建','api_relation_create','POST','/new/api_relation','2021-03-22 09:29:01','2021-03-22 09:29:01'),(134,'api_relation 删除','api_relation_delete','POST','/delete/api_relation','2021-03-22 09:29:01','2021-03-22 09:29:01'),(135,'api_relation 导出','api_relation_export','POST','/export/api_relation','2021-03-22 09:29:01','2021-03-22 09:29:01'),(136,'api_fuzzing_data 查询','api_fuzzing_data_query','GET','/info/api_fuzzing_data','2021-03-22 09:34:15','2021-03-22 09:34:15'),(137,'api_fuzzing_data 编辑页显示','api_fuzzing_data_show_edit','GET','/info/api_fuzzing_data/edit','2021-03-22 09:34:15','2021-03-22 09:34:15'),(138,'api_fuzzing_data 新建记录页显示','api_fuzzing_data_show_create','GET','/info/api_fuzzing_data/new','2021-03-22 09:34:15','2021-03-22 09:34:15'),(139,'api_fuzzing_data 编辑','api_fuzzing_data_edit','POST','/edit/api_fuzzing_data','2021-03-22 09:34:15','2021-03-22 09:34:15'),(140,'api_fuzzing_data 新建','api_fuzzing_data_create','POST','/new/api_fuzzing_data','2021-03-22 09:34:15','2021-03-22 09:34:15'),(141,'api_fuzzing_data 删除','api_fuzzing_data_delete','POST','/delete/api_fuzzing_data','2021-03-22 09:34:15','2021-03-22 09:34:15'),(142,'api_fuzzing_data 导出','api_fuzzing_data_export','POST','/export/api_fuzzing_data','2021-03-22 09:34:15','2021-03-22 09:34:15'),(143,'api_id_count 查询','api_id_count_query','GET','/info/api_id_count','2021-03-22 09:41:12','2021-03-22 09:41:12'),(144,'api_id_count 编辑页显示','api_id_count_show_edit','GET','/info/api_id_count/edit','2021-03-22 09:41:12','2021-03-22 09:41:12'),(145,'api_id_count 新建记录页显示','api_id_count_show_create','GET','/info/api_id_count/new','2021-03-22 09:41:12','2021-03-22 09:41:12'),(146,'api_id_count 编辑','api_id_count_edit','POST','/edit/api_id_count','2021-03-22 09:41:12','2021-03-22 09:41:12'),(147,'api_id_count 新建','api_id_count_create','POST','/new/api_id_count','2021-03-22 09:41:12','2021-03-22 09:41:12'),(148,'api_id_count 删除','api_id_count_delete','POST','/delete/api_id_count','2021-03-22 09:41:12','2021-03-22 09:41:12'),(149,'api_id_count 导出','api_id_count_export','POST','/export/api_id_count','2021-03-22 09:41:12','2021-03-22 09:41:12'),(150,'product_count 查询','product_count_query','GET','/info/product_count','2021-03-22 09:43:58','2021-03-22 09:43:58'),(151,'product_count 编辑页显示','product_count_show_edit','GET','/info/product_count/edit','2021-03-22 09:43:58','2021-03-22 09:43:58'),(152,'product_count 新建记录页显示','product_count_show_create','GET','/info/product_count/new','2021-03-22 09:43:58','2021-03-22 09:43:58'),(153,'product_count 编辑','product_count_edit','POST','/edit/product_count','2021-03-22 09:43:58','2021-03-22 09:43:58'),(154,'product_count 新建','product_count_create','POST','/new/product_count','2021-03-22 09:43:58','2021-03-22 09:43:58'),(155,'product_count 删除','product_count_delete','POST','/delete/product_count','2021-03-22 09:43:58','2021-03-22 09:43:58'),(156,'product_count 导出','product_count_export','POST','/export/product_count','2021-03-22 09:43:58','2021-03-22 09:43:58'),(157,'parameter_definition 查询','parameter_definition_query','GET','/info/parameter_definition','2021-03-22 09:45:34','2021-03-22 09:45:34'),(158,'parameter_definition 编辑页显示','parameter_definition_show_edit','GET','/info/parameter_definition/edit','2021-03-22 09:45:34','2021-03-22 09:45:34'),(159,'parameter_definition 新建记录页显示','parameter_definition_show_create','GET','/info/parameter_definition/new','2021-03-22 09:45:34','2021-03-22 09:45:34'),(160,'parameter_definition 编辑','parameter_definition_edit','POST','/edit/parameter_definition','2021-03-22 09:45:34','2021-03-22 09:45:34'),(161,'parameter_definition 新建','parameter_definition_create','POST','/new/parameter_definition','2021-03-22 09:45:34','2021-03-22 09:45:34'),(162,'parameter_definition 删除','parameter_definition_delete','POST','/delete/parameter_definition','2021-03-22 09:45:34','2021-03-22 09:45:34'),(163,'parameter_definition 导出','parameter_definition_export','POST','/export/parameter_definition','2021-03-22 09:45:34','2021-03-22 09:45:34'),(164,'fuzzing_definition 查询','fuzzing_definition_query','GET','/info/fuzzing_definition','2021-03-24 02:12:43','2021-03-24 02:12:43'),(165,'fuzzing_definition 编辑页显示','fuzzing_definition_show_edit','GET','/info/fuzzing_definition/edit','2021-03-24 02:12:43','2021-03-24 02:12:43'),(166,'fuzzing_definition 新建记录页显示','fuzzing_definition_show_create','GET','/info/fuzzing_definition/new','2021-03-24 02:12:43','2021-03-24 02:12:43'),(167,'fuzzing_definition 编辑','fuzzing_definition_edit','POST','/edit/fuzzing_definition','2021-03-24 02:12:43','2021-03-24 02:12:43'),(168,'fuzzing_definition 新建','fuzzing_definition_create','POST','/new/fuzzing_definition','2021-03-24 02:12:43','2021-03-24 02:12:43'),(169,'fuzzing_definition 删除','fuzzing_definition_delete','POST','/delete/fuzzing_definition','2021-03-24 02:12:43','2021-03-24 02:12:43'),(170,'fuzzing_definition 导出','fuzzing_definition_export','POST','/export/fuzzing_definition','2021-03-24 02:12:43','2021-03-24 02:12:43'),(171,'scene_data 查询','scene_data_query','GET','/info/scene_data','2021-06-03 03:34:06','2021-06-03 03:34:06'),(172,'scene_data 编辑页显示','scene_data_show_edit','GET','/info/scene_data/edit','2021-06-03 03:34:06','2021-06-03 03:34:06'),(173,'scene_data 新建记录页显示','scene_data_show_create','GET','/info/scene_data/new','2021-06-03 03:34:06','2021-06-03 03:34:06'),(174,'scene_data 编辑','scene_data_edit','POST','/edit/scene_data','2021-06-03 03:34:06','2021-06-03 03:34:06'),(175,'scene_data 新建','scene_data_create','POST','/new/scene_data','2021-06-03 03:34:06','2021-06-03 03:34:06'),(176,'scene_data 删除','scene_data_delete','POST','/delete/scene_data','2021-06-03 03:34:06','2021-06-03 03:34:06'),(177,'scene_data 导出','scene_data_export','POST','/export/scene_data','2021-06-03 03:34:06','2021-06-03 03:34:06'),(178,'playbook 查询','playbook_query','GET','/info/playbook','2021-06-03 03:34:54','2021-06-03 03:34:54'),(179,'playbook 编辑页显示','playbook_show_edit','GET','/info/playbook/edit','2021-06-03 03:34:54','2021-06-03 03:34:54'),(180,'playbook 新建记录页显示','playbook_show_create','GET','/info/playbook/new','2021-06-03 03:34:54','2021-06-03 03:34:54'),(181,'playbook 编辑','playbook_edit','POST','/edit/playbook','2021-06-03 03:34:54','2021-06-03 03:34:54'),(182,'playbook 新建','playbook_create','POST','/new/playbook','2021-06-03 03:34:54','2021-06-03 03:34:54'),(183,'playbook 删除','playbook_delete','POST','/delete/playbook','2021-06-03 03:34:54','2021-06-03 03:34:54'),(184,'playbook 导出','playbook_export','POST','/export/playbook','2021-06-03 03:34:54','2021-06-03 03:34:54'),(185,'apiManage','apiManage','GET','/likePostman','2022-03-07 02:31:22','2022-03-07 02:31:22'),(186,'sys_parameter 查询','sys_parameter_query','GET','/info/sys_parameter','2021-08-05 18:39:20','2021-08-05 18:39:20'),(187,'sys_parameter 编辑页显示','sys_parameter_show_edit','GET','/info/sys_parameter/edit','2021-08-05 18:39:20','2021-08-05 18:39:20'),(188,'sys_parameter 新建记录页显示','sys_parameter_show_create','GET','/info/sys_parameter/new','2021-08-05 18:39:20','2021-08-05 18:39:20'),(189,'sys_parameter 编辑','sys_parameter_edit','POST','/edit/sys_parameter','2021-08-05 18:39:20','2021-08-05 18:39:20'),(190,'sys_parameter 新建','sys_parameter_create','POST','/new/sys_parameter','2021-08-05 18:39:20','2021-08-05 18:39:20'),(191,'sys_parameter 删除','sys_parameter_delete','POST','/delete/sys_parameter','2021-08-05 18:39:20','2021-08-05 18:39:20'),(192,'sys_parameter 导出','sys_parameter_export','POST','/export/sys_parameter','2021-08-05 18:39:20','2021-08-05 18:39:20'),(193,'app_api_changelog 查询','app_api_changelog_query','GET','/info/app_api_changelog','2022-03-07 11:48:47','2022-03-07 11:48:47'),(194,'app_api_changelog 编辑页显示','app_api_changelog_show_edit','GET','/info/app_api_changelog/edit','2022-03-07 11:48:47','2022-03-07 11:48:47'),(195,'app_api_changelog 新建记录页显示','app_api_changelog_show_create','GET','/info/app_api_changelog/new','2022-03-07 11:48:47','2022-03-07 11:48:47'),(196,'app_api_changelog 编辑','app_api_changelog_edit','POST','/edit/app_api_changelog','2022-03-07 11:48:47','2022-03-07 11:48:47'),(197,'app_api_changelog 新建','app_api_changelog_create','POST','/new/app_api_changelog','2022-03-07 11:48:47','2022-03-07 11:48:47'),(198,'app_api_changelog 删除','app_api_changelog_delete','POST','/delete/app_api_changelog','2022-03-07 11:48:47','2022-03-07 11:48:47'),(199,'app_api_changelog 导出','app_api_changelog_export','POST','/export/app_api_changelog','2022-03-07 11:48:47','2022-03-07 11:48:47'),(200,'app_dashboard','app_dashboard','GET','/app_dashboard','2022-03-08 11:35:21','2022-03-09 03:52:29'),(213,'scene_data_test_history 查询','scene_data_test_history_query','GET','/info/scene_data_test_history','2022-03-09 06:21:00','2022-03-09 06:21:00'),(214,'scene_data_test_history 编辑页显示','scene_data_test_history_show_edit','GET','/info/scene_data_test_history/edit','2022-03-09 06:21:00','2022-03-09 06:21:00'),(215,'scene_data_test_history 新建记录页显示','scene_data_test_history_show_create','GET','/info/scene_data_test_history/new','2022-03-09 06:21:00','2022-03-09 06:21:00'),(216,'scene_data_test_history 编辑','scene_data_test_history_edit','POST','/edit/scene_data_test_history','2022-03-09 06:21:00','2022-03-09 06:21:00'),(217,'scene_data_test_history 新建','scene_data_test_history_create','POST','/new/scene_data_test_history','2022-03-09 06:21:00','2022-03-09 06:21:00'),(218,'scene_data_test_history 删除','scene_data_test_history_delete','POST','/delete/scene_data_test_history','2022-03-09 06:21:00','2022-03-09 06:21:00'),(219,'scene_data_test_history 导出','scene_data_test_history_export','POST','/export/scene_data_test_history','2022-03-09 06:21:00','2022-03-09 06:21:00'),(220,'scene_test_history 查询','scene_test_history_query','GET','/info/scene_test_history','2022-03-09 06:21:10','2022-03-09 06:21:10'),(221,'scene_test_history 编辑页显示','scene_test_history_show_edit','GET','/info/scene_test_history/edit','2022-03-09 06:21:10','2022-03-09 06:21:10'),(222,'scene_test_history 新建记录页显示','scene_test_history_show_create','GET','/info/scene_test_history/new','2022-03-09 06:21:10','2022-03-09 06:21:10'),(223,'scene_test_history 编辑','scene_test_history_edit','POST','/edit/scene_test_history','2022-03-09 06:21:10','2022-03-09 06:21:10'),(224,'scene_test_history 新建','scene_test_history_create','POST','/new/scene_test_history','2022-03-09 06:21:10','2022-03-09 06:21:10'),(225,'scene_test_history 删除','scene_test_history_delete','POST','/delete/scene_test_history','2022-03-09 06:21:10','2022-03-09 06:21:10'),(226,'scene_test_history 导出','scene_test_history_export','POST','/export/scene_test_history','2022-03-09 06:21:10','2022-03-09 06:21:10'),(227,'history_data_preview ','history_data_preview','GET','/fm/history/preview','2022-03-09 06:28:26','2022-03-09 06:28:26'),(228,'scene_data_preview','scene_data_preview','GET','/fm/scene/preview','2022-03-09 06:30:30','2022-03-09 06:30:30'),(229,'assert_template 查询','assert_template_query','GET','/info/assert_template','2023-12-11 07:20:38','2023-12-11 07:20:38'),(230,'assert_template 编辑页显示','assert_template_show_edit','GET','/info/assert_template/edit','2023-12-11 07:20:38','2023-12-11 07:20:38'),(231,'assert_template 新建记录页显示','assert_template_show_create','GET','/info/assert_template/new','2023-12-11 07:20:38','2023-12-11 07:20:38'),(232,'assert_template 编辑','assert_template_edit','POST','/edit/assert_template','2023-12-11 07:20:38','2023-12-11 07:20:38'),(233,'assert_template 新建','assert_template_create','POST','/new/assert_template','2023-12-11 07:20:38','2023-12-11 07:20:38'),(234,'assert_template 删除','assert_template_delete','POST','/delete/assert_template','2023-12-11 07:20:38','2023-12-11 07:20:38'),(235,'assert_template 导出','assert_template_export','POST','/export/assert_template','2023-12-11 07:20:38','2023-12-11 07:20:38');
/*!40000 ALTER TABLE `goadmin_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goadmin_site`
--

LOCK TABLES `goadmin_site` WRITE;
/*!40000 ALTER TABLE `goadmin_site` DISABLE KEYS */;
INSERT INTO `goadmin_site` VALUES (1,'theme','sword',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(2,'hide_app_info_entrance','true',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(3,'info_log_path','./mgmt/log/info.log',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(4,'info_log_off','false',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(5,'mini_logo','        盾测\r\n    ',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(6,'sql_log','false',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(7,'login_logo','       盾测-自动化\r\n    ',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(8,'hide_plugin_entrance','true',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(9,'custom_500_html','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(10,'login_title','盾测-自动化',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(11,'domain','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(12,'store','{\"path\":\"./uploads\",\"prefix\":\"uploads\"}',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(13,'logger_rotate_max_age','30',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(14,'logger_encoder_time_key','ts',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(15,'asset_url','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(16,'custom_404_html','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(17,'animation_duration','0.00',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(18,'allow_del_operation_log','true',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(19,'logger_rotate_compress','false',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(20,'logger_encoder_encoding','console',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(21,'custom_head_html','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(22,'animation_type','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(23,'hide_tool_entrance','true',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(24,'title','盾测',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(25,'logger_level','-1',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(26,'custom_403_html','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(27,'operation_log_off','false',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(28,'logger_rotate_max_size','10',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(29,'auth_user_table','goadmin_users',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(30,'login_url','/login',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(31,'error_log_path','./mgmt/log/error.log',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(32,'logger_rotate_max_backups','5',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(33,'go_mod_file_path','./mgmt/commongo.mod',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(34,'animation_delay','0.00',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(35,'logger_encoder_time','iso8601',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(36,'custom_foot_html','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(37,'asset_root_path','./public/',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(38,'logger_encoder_message_key','msg',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(39,'databases','{\"default\":{\"host\":\"127.0.0.1\",\"port\":\"3306\",\"user\":\"root\",\"pwd\":\"zaq1@WSX\",\"name\":\"data4perf\",\"max_idle_con\":5,\"max_open_con\":10,\"driver\":\"mysql\"}}',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(40,'url_prefix','api/v1',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(41,'logger_encoder_level_key','level',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(42,'logger_encoder_caller_key','caller',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(43,'logger_encoder_stacktrace_key','stacktrace',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(44,'hide_config_center_entrance','true',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(45,'access_log_off','false',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(46,'access_assets_log_off','false',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(47,'logger_encoder_caller','short',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(48,'session_life_time','7200',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(49,'language','zh',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(50,'logo','        盾测-自动化\r\n    ',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(51,'logger_encoder_name_key','logger',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(52,'logger_encoder_duration','string',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(53,'error_log_off','false',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(54,'app_id','j90eXvI3x1ye',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(55,'extra','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(56,'index_url','/',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(57,'env','local',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(58,'color_scheme','skin-black',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(59,'file_upload_engine','{\"name\":\"local\"}',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(60,'bootstrap_file_path','./mgmt/common/bootstrap.go',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(61,'footer_info','',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(62,'site_off','false',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(63,'debug','true',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(64,'access_log_path','./mgmt/log/access.log',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(65,'logger_encoder_level','capitalColor',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(66,'no_limit_login_ip','false',NULL,1,'2020-11-23 02:04:52','2020-11-23 02:04:52'),(67,'filemanager_connection','default',NULL,0,'2020-12-07 02:36:00','2020-12-07 02:36:00'),(68,'hide_visitor_user_center_entrance','false',NULL,1,'2021-03-22 12:48:28','2021-03-22 12:48:28'),(69,'prohibit_config_modification','false',NULL,1,'2021-03-22 12:48:28','2021-03-22 12:48:28'),(70,'exclude_theme_components','null',NULL,1,'2021-03-22 12:48:28','2021-03-22 12:48:28'),(71,'open_admin_api','false',NULL,1,'2021-03-22 12:48:28','2021-03-22 12:48:28'),(72,'librarian_build_menu_def','2564,2565,2566,2567,2568,2569,2570,2571,2572,2573,2574,2575,2576,2577',NULL,0,'2021-03-23 05:41:47','2021-03-23 05:41:47'),(73,'librarian_build_menu_def_nav','41ef6298be0e85ff535f4814e18245d5',NULL,0,'2021-03-23 05:41:47','2021-03-23 05:41:47'),(74,'librarian_roots','{\"mk\":{\"Path\":\"./mgmt/data4perf/doc/file\",\"Title\":\"mk\"},\"key\": \"系统说明\"}',NULL,0,'2021-03-23 05:43:38','2021-03-23 05:43:38'),(75,'librarian_theme','github',NULL,0,'2021-03-23 05:43:38','2021-03-23 05:43:38');
/*!40000 ALTER TABLE `goadmin_site` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-01-05 20:54:21
-- MySQL dump 10.14  Distrib 5.5.60-MariaDB, for Linux (x86_64)
--
-- Host: 10.57.242.40    Database: data4test
-- ------------------------------------------------------
-- Server version	5.7.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `goadmin_role_menu`
--

LOCK TABLES `goadmin_role_menu` WRITE;
/*!40000 ALTER TABLE `goadmin_role_menu` DISABLE KEYS */;
INSERT INTO `goadmin_role_menu` VALUES (1,1,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(1,8,'2020-11-23 02:57:07','2020-11-23 02:57:07'),(2,8,'2020-11-23 02:57:07','2020-11-23 02:57:07'),(1,18,'2020-11-23 03:06:44','2020-11-23 03:06:44'),(2,18,'2020-11-23 03:06:44','2020-11-23 03:06:44'),(1,34,'2020-12-07 03:06:28','2020-12-07 03:06:28'),(2,34,'2020-12-07 03:06:28','2020-12-07 03:06:28'),(1,37,'2020-12-07 03:08:41','2020-12-07 03:08:41'),(2,37,'2020-12-07 03:08:41','2020-12-07 03:08:41'),(1,38,'2020-12-07 03:09:05','2020-12-07 03:09:05'),(2,38,'2020-12-07 03:09:05','2020-12-07 03:09:05'),(1,42,'2020-12-07 03:24:04','2020-12-07 03:24:04'),(2,42,'2020-12-07 03:24:04','2020-12-07 03:24:04'),(1,43,'2020-12-08 08:20:47','2020-12-08 08:20:47'),(2,43,'2020-12-08 08:20:47','2020-12-08 08:20:47'),(1,47,'2020-12-10 03:10:35','2020-12-10 03:10:35'),(2,47,'2020-12-10 03:10:35','2020-12-10 03:10:35'),(1,45,'2020-12-10 03:16:07','2020-12-10 03:16:07'),(2,45,'2020-12-10 03:16:07','2020-12-10 03:16:07'),(1,39,'2020-12-10 03:17:02','2020-12-10 03:17:02'),(2,39,'2020-12-10 03:17:02','2020-12-10 03:17:02'),(1,44,'2020-12-10 03:17:15','2020-12-10 03:17:15'),(2,44,'2020-12-10 03:17:15','2020-12-10 03:17:15'),(1,7,'2020-12-28 14:16:17','2020-12-28 14:16:17'),(2,7,'2020-12-28 14:16:17','2020-12-28 14:16:17'),(1,50,'2021-03-18 14:53:17','2021-03-18 14:53:17'),(1,22,'2021-03-22 09:47:56','2021-03-22 09:47:56'),(2,22,'2021-03-22 09:47:56','2021-03-22 09:47:56'),(1,25,'2021-03-22 09:48:08','2021-03-22 09:48:08'),(2,25,'2021-03-22 09:48:08','2021-03-22 09:48:08'),(1,51,'2021-03-22 11:26:03','2021-03-22 11:26:03'),(2,51,'2021-03-22 11:26:03','2021-03-22 11:26:03'),(1,17,'2021-03-30 03:02:07','2021-03-30 03:02:07'),(2,17,'2021-03-30 03:02:07','2021-03-30 03:02:07'),(1,9,'2021-07-29 06:04:57','2021-07-29 06:04:57'),(2,9,'2021-07-29 06:04:57','2021-07-29 06:04:57'),(3,9,'2022-03-09 07:19:56','2022-03-09 07:19:56'),(3,8,'2022-03-09 07:20:47','2022-03-09 07:20:47'),(5,236,'2022-03-09 07:30:24','2022-03-09 07:30:24'),(1,46,'2023-08-11 10:00:56','2023-08-11 10:00:56'),(2,46,'2023-08-11 10:00:56','2023-08-11 10:00:56'),(1,48,'2023-08-11 10:01:05','2023-08-11 10:01:05'),(2,48,'2023-08-11 10:01:05','2023-08-11 10:01:05'),(1,41,'2023-08-11 10:05:06','2023-08-11 10:05:06'),(2,41,'2023-08-11 10:05:06','2023-08-11 10:05:06'),(1,40,'2023-08-11 10:05:23','2023-08-11 10:05:23'),(2,40,'2023-08-11 10:05:23','2023-08-11 10:05:23'),(3,2564,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2565,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2566,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2567,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2568,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2569,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2570,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2571,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2572,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2573,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2574,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2575,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2576,'2024-01-05 12:45:11','2024-01-05 12:45:11'),(3,2577,'2024-01-05 12:45:11','2024-01-05 12:45:11');
/*!40000 ALTER TABLE `goadmin_role_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goadmin_role_permissions`
--

LOCK TABLES `goadmin_role_permissions` WRITE;
/*!40000 ALTER TABLE `goadmin_role_permissions` DISABLE KEYS */;
INSERT INTO `goadmin_role_permissions` VALUES (1,1,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(1,2,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(2,2,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,3,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,4,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,5,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,6,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,7,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,8,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,9,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,10,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,11,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,12,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,13,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,14,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,15,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,16,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,17,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,18,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,19,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,20,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,21,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,22,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,23,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,24,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,25,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,26,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,27,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,28,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,29,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,30,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,31,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,32,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,33,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,34,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,35,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,36,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,37,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,38,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,39,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,40,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,41,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,42,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,43,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,44,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,45,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,46,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,47,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,48,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,49,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,50,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,51,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,52,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,53,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,54,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,55,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,56,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,57,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,58,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,59,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,60,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,61,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,62,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,63,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,64,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,65,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,66,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,67,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,68,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,69,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,70,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,71,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,72,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,73,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,74,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,75,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,76,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,77,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,78,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,79,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,80,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,81,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,82,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,83,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,84,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,85,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,86,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,87,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,88,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,89,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,90,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,91,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,92,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,93,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,94,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,95,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,96,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,97,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,98,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,99,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,100,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,101,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,102,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,103,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,104,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,105,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,106,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,107,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,108,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,109,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,110,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,111,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,112,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,113,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,114,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,115,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,116,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,117,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,118,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,119,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,120,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,121,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,122,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,123,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,124,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,125,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,126,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,127,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,128,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,129,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,130,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,131,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,132,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,133,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,134,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,135,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,136,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,137,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,138,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,139,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,140,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,141,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,142,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,143,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,144,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,145,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,146,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,147,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,148,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,149,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,150,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,151,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,152,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,153,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,154,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,155,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,156,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,157,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,158,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,159,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,160,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,161,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,162,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,163,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,164,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,165,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,166,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,167,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,168,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,169,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,170,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,171,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,172,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,173,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,174,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,175,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,176,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,177,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,178,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,179,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,180,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,181,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,182,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,183,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,184,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,185,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,186,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,187,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,188,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,189,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,190,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,191,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,192,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,193,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,194,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,195,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,196,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,197,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,198,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,199,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,200,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,213,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,214,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,215,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,216,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,217,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,218,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,219,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,220,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,221,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,222,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,223,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,224,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,225,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,226,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,227,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(2,228,'2023-06-16 02:26:56','2023-06-16 02:26:56'),(3,3,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,115,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,122,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,123,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,171,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,172,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,200,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,213,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,214,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,220,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,221,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,227,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(3,228,'2022-03-09 07:27:55','2022-03-09 07:27:55'),(5,185,'2022-03-09 07:27:47','2022-03-09 07:27:47');
/*!40000 ALTER TABLE `goadmin_role_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goadmin_role_users`
--

LOCK TABLES `goadmin_role_users` WRITE;
/*!40000 ALTER TABLE `goadmin_role_users` DISABLE KEYS */;
INSERT INTO `goadmin_role_users` VALUES (1,1,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(2,2,'2019-09-09 16:00:00','2019-09-09 16:00:00');
/*!40000 ALTER TABLE `goadmin_role_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goadmin_roles`
--

LOCK TABLES `goadmin_roles` WRITE;
/*!40000 ALTER TABLE `goadmin_roles` DISABLE KEYS */;
INSERT INTO `goadmin_roles` VALUES (1,'Administrator','administrator','2019-09-09 16:00:00','2019-09-09 16:00:00'),(2,'Operator','operator','2019-09-09 16:00:00','2023-06-16 02:26:56'),(3,'Cicd','cicd','2022-03-07 02:30:02','2022-03-09 07:27:55'),(5,'ApiManage','apiManage','2022-03-09 07:27:36','2022-03-09 07:27:47');
/*!40000 ALTER TABLE `goadmin_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goadmin_users`
--

LOCK TABLES `goadmin_users` WRITE;
/*!40000 ALTER TABLE `goadmin_users` DISABLE KEYS */;
INSERT INTO `goadmin_users` VALUES (1,'admin','$2a$10$2JwnVImEmgieSFBYz3uKhu2zI86ytpWiaaniFGdFu72tg9/mOsnaq','admin','','tlNcBVK9AvfYH7WEnwB1RKvocJu8FfRy4um3DJtwdHuJy0dwFsLOgAc0xUfh','2019-09-09 16:00:00','2019-09-09 16:00:00'),(2,'operator','$2a$10$Y8BSfJuwRBZ9pxgzaWpCnub0eja4XE93zbkzpep7GawO8BCJ3fK2C','Operator','',NULL,'2019-09-09 16:00:00','2019-09-09 16:00:00');
/*!40000 ALTER TABLE `goadmin_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `goadmin_user_permissions`
--

LOCK TABLES `goadmin_user_permissions` WRITE;
/*!40000 ALTER TABLE `goadmin_user_permissions` DISABLE KEYS */;
INSERT INTO `goadmin_user_permissions` VALUES (1,1,'2019-09-09 16:00:00','2019-09-09 16:00:00'),(2,2,'2019-09-09 16:00:00','2019-09-09 16:00:00');
/*!40000 ALTER TABLE `goadmin_user_permissions` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-01-05 20:54:41
-- MySQL dump 10.14  Distrib 5.5.60-MariaDB, for Linux (x86_64)
--
-- Host: 10.57.242.40    Database: data4test
-- ------------------------------------------------------
-- Server version	5.7.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `assert_template`
--

LOCK TABLES `assert_template` WRITE;
/*!40000 ALTER TABLE `assert_template` DISABLE KEYS */;
INSERT INTO `assert_template` VALUES (1,'successTemplate','{\"ch\": \"成功|重复|已存在|已经存在\", \"en\": \"success|Success|exist|duplicate\"}','数据文件断言值中中以 \'{successTemplate}\'关联，进行断言时会根据请求语种进行结果判断','admin','2023-12-11 08:08:37','2023-12-11 08:10:59',NULL);
/*!40000 ALTER TABLE `assert_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `env_config`
--

LOCK TABLES `env_config` WRITE;
/*!40000 ALTER TABLE `env_config` DISABLE KEYS */;
INSERT INTO `env_config` VALUES (1,'示例产品','exampleApp','X.X.X.X:8088','http','/prefix','no',1,'{\"Content-Type\":\"application/x-www-form-urlencoded\",\"Cookie\":\"xxx\",\"Referer\":\"http://x.x.x.x:80xx\",\"X-Cf-Random\":\"xxx\"}','custom','http://x.x.x.x:80xx/api/v2/api-docs?group=group1','','2023-12-13 01:59:03','2023-12-13 02:17:14',NULL);
/*!40000 ALTER TABLE `env_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `sys_parameter`
--

LOCK TABLES `sys_parameter` WRITE;
/*!40000 ALTER TABLE `sys_parameter` DISABLE KEYS */;
INSERT INTO `sys_parameter` VALUES (1,'fileName','file, gradeModel, model, pickleFiles, templateFile','当header头中Content-Type值multipart/form-data时，请求参数名称为值列表中的名称时，请求时以文件模式加载','2021-08-06 03:56:05','2021-08-06 09:14:51',NULL);
/*!40000 ALTER TABLE `sys_parameter` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `scene_data`
--

LOCK TABLES `scene_data` WRITE;
/*!40000 ALTER TABLE `scene_data` DISABLE KEYS */;
INSERT INTO `scene_data` VALUES (1,'示例-用户管理-新建用户','post_/path','exampleApp','示例-用户管理-新建用户.yml','---\r\n# 用例信息\r\nname: 示例-用户管理-新建用户 # 数据用例名称，e.g.: 类型-模块-用例， 类型：功能/性能/异常/内置/……， 模块：用户管理/规则管理/……\r\napi_id: post_/path        # 用例ID, method_path组合，后续做数据联动使用，数据统计使用\r\nversion: 1.0              # 数据用例版本，后续可以进行数据升级\r\nis_run_pre_apis: \"no\"     # 是否跑前置用例，选项：yes / no,  默认 no， 功能未开发\r\nis_run_post_apis: \"no\"    # 是否跑后置用例，选项：yes / no,  默认 no， 功能未开发\r\nis_parallel: \"no\"         # 是否并行跑数据，选项：yes / no,  默认 no，\r\nis_use_env_config: \"yes\"  # 是否使用公共环境，选项：yes / no,  默认 yes\r\n\r\n# 环境信息\r\nenv:\r\n  protocol: http        # http 或 https，请求协议\r\n  host: X.X.X.X:8089    # 环境IP 或 环境域名 或 环境IP:端口\r\n  prepath: /prefix      # 路由前缀，公共部分可以抽出来\r\n\r\n# API 基本信息\r\napi:\r\n  description: 新建用户   # API用途\r\n  module: 用户管理        # API所属模块\r\n  app: appName           # API所属应用\r\n  method: post           # （注意：保证正确） API请求方法\r\n  path: /path            # （注意：保证正确）API请求路由，路由前缀抽离到prepath下时或公共环境中已定义prepath时，这里无需再写路由前缀\r\n  pre_apis: []           # 调试时，依赖前置用例时，可以把关联前置文件写上，功能未充分验证\r\n  param_apis: []         # 调试时，依赖其他用例的参数时，可以把关联文件写上，功能未充分验证\r\n  post_apis: []          # 调试时，测试跑完后需要跑的用例，可以把关联文件写上，功能未充分验证\r\n\r\n# 定义单值参数，如果is_use_env_config值为no, 需要定义此处的 header\r\nsingle:\r\n  header:\r\n    Content-Type: multipart/form-data   # 如果api为导入文件功能，需要把Content-Type定义为multipart/form-data进行公用环境值的覆盖，优化级：数据文件>应用配置>产品配置\r\n  query: {}                             # GET请求时，请求参数定义，只定义一个值，共用的参数放在这里，无需反复定义\r\n  path: {}                              # PATH 变量参数定义，只定义一个值\r\n  body: \r\n    uuid: \'{self}\'                      # \'{self}\'代表该值从依赖用例的输出参数中获取\r\n    condition: \'{\"children\":[{\"logicOperator\":\"&&\",\"property\":\"watchlist/customList\",\"operator\":\"==\",\"value\":\"1\",\"type\":\"alias\",\"description\":\"\",\"propertyDataType\":\"\",\"children\":[],\"describe\":\"是否命中名单\",\"params\":[{\"name\":\"calcField\",\"type\":\"string\",\"value\":\"S_DC_VS_NAME\"},{\"name\":\"definitionList\",\"type\":\"string\",\"value\":\"{nameList}\"},{\"name\":\"conditionValue\",\"type\":\"int\",\"value\":\"1\"}]}],\"logicOperator\":\"&&\"}\'  # {nameList} 代表字符串里有需替换的变更，nameList为 ouput 中输出的参数名字，在前置的用例中有定义同名变量，即会替换\r\n\r\n# 定义多值参数\r\nmulti:\r\n  query: {}                   # GET请求时，请求参数定义，定义的值为列表\r\n  path: {}                    # PATH 变更参数定义，定义的值为列表\r\n  body:\r\n    description:              # 定义多值时，取各项定义的个数最少的数据，一一对应\r\n    - \'{Rune(128)}\'    # 获取设置长度的汉字\r\n    - \'{Str(64)}\'      # 获取设置长度的字符串\r\n    - \'{Int(10,100)}\'  # 获取设置范围内的整数\r\n    displayName:\r\n    - \'{Date(-2)}\'      # 获取两天前的日期\r\n    - \'{Date(2)}\'       # 获取两天后的日期\r\n    - \'{Timestamp(-2)}\' # 获取两天前的时间戳\r\n    name:\r\n    - \'{IDNo}\'          # 获取身份证字符串\r\n    - \'{Name}\'          # 获取名字字符串\r\n    - \'{Address}\'       # 获取地址字符串\r\n    - \'{BankNo}\'        # 获取银行卡号字符串\r\n\r\n# 断言，数据校验，根据需要写不同类型的断言，不写断言，只要返回为200，即算 PASS\r\nassert:\r\n- type: equal   # 验证code的值等于200\r\n  source: code    # 返回的json信息，取key为code的值\r\n  value: 200 \r\n- type: \"!=\" # 验证code的值不等于200\r\n  source: code    # 返回的json信息，取key为code的值\r\n  value: 200 \r\n- type: \">=\"    # 验证source字段大于等于1\r\n  source: data-total     # 返回的json信息，data字典-取出productDesc的值\r\n  value: 1\r\n- type:  contain\r\n  source: data-contents*productDesc  # 返回的json信息，data字典-content字典-字典列表，取出productDesc的值, 并校验是否包含 value中的值\r\n  value: 待删除的产品描述\r\n- type: \"!in\"   # 验证取到的productName的值包含删除\r\n  source: data-contents*productName  # 返回的json信息，data字典-content字典-字典列表，取出productName的值, 不包含 value 中的值\r\n  value: 删除\r\n- type: not_contain   # 验证取到的productName的值不包含删除\r\n  source: data-contents*productName  # 返回的json信息，data字典-content字典-字典列表，取出productName的值\r\n  value: 产品\r\n- type: re\r\n  source: message\r\n  value: 成功|重复|已存在\r\n- type: output  # 从返回的json 信息取取出 uuid 的值，并命名为uuid\r\n  source: data-contents*uuid\r\n  value: uuid\r\n- type: output  # 从返回的json信息取出uuid的值，并重命名为ProductUuid\r\n  source: data-contents*uuid\r\n  value: ProductUuid\r\n# 数据执行后的动作\r\naction:\r\n  - type: sleep\r\n    value: 1    // 表示等待1秒种，时间可根据需要自动设置，单位为秒\r\n  - type: create_csv\r\n    value: name:number    // 生成文件名:生成的数据条数，默认生成10条\r\n  - type: create_xls\r\n    value: name:number    // 生成文件名:生成的数据条数, 默认生成10条\r\n# 输出其他接口需要的依赖数据, 由断言中类型为 ouput 定义，自动生成, 定义为\'{self}\'或 \'{uuid}\' 从此处取值\r\noutput:\r\n  uuid:\r\n  - XXX\r\n  - XXX\r\n\r\n# 测试结果：pass, fail, untest, 自动生成，断言全部符合要求设为pass, 请求若返回非200，直接置为 fail, 如果执行次数测试为0，测置为 untest\r\n# 保留最新测试结果\r\ntest_result:\r\n- pass\r\n- fail\r\n- untest\r\n\r\n# 请求 URL，自动生成， 保留最新测试结果\r\nurls:\r\n- http://X.X.X.X:8089/creditApi/loanProduct/pageProductList\r\n\r\n# 请求数据，body, query, 自动生成, 保留最新测试结果\r\nrequests:\r\n  - \'{\"curPage\":\"1\",\"endEntryTime\":\"1627095420000\",\"pageSize\":\"10\",\"searchOption\":\"{}\"startEntryTime\":\"1626749820000\",\"timeType\":\"1\"}\'\r\n\r\n# 返回信息, 自动生成， 保留最新测试结果\r\nresponse:\r\n- \"response1\"\r\n- \"response2\"',1,'','','','admin','2023-12-13 02:20:32',NULL,NULL);
/*!40000 ALTER TABLE `scene_data` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `product`
--

LOCK TABLES `product` WRITE;
/*!40000 ALTER TABLE `product` DISABLE KEYS */;
INSERT INTO `product` VALUES (1,'示例产品','X.X.X.X:8088','http','no',1,'{\r\n    \"Content-Type\": \"application/x-www-form-urlencoded\",\r\n    \"Cookie\": \"XXXX\",\r\n    \"Referer\": \"http://X.X.X.X:8088/\",\r\n    \"X-Cf-Random\": \"xxx\"\r\n}','custom','',2,'2023-12-13 01:57:50','{\"name1\": \"value1\", \"name2\": \"value2\", \"name3\": \"value3\"}','','2023-12-13 02:16:57',NULL);
/*!40000 ALTER TABLE `product` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `playbook`
--

LOCK TABLES `playbook` WRITE;
/*!40000 ALTER TABLE `playbook` DISABLE KEYS */;
INSERT INTO `playbook` VALUES (1,'示例场景','1','<a href=\"/admin/fm/data/preview?path=/示例-用户管理-新建用户.yml\">示例-用户管理-新建用户.yml</a>','','1',1,1,'','',NULL,'admin','示例产品','2023-12-13 02:23:25',NULL,NULL);
/*!40000 ALTER TABLE `playbook` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-01-05 20:55:11
-- MySQL dump 10.14  Distrib 5.5.60-MariaDB, for Linux (x86_64)
--
-- Host: 10.57.242.40    Database: data4test
-- ------------------------------------------------------
-- Server version	5.7.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `schedule`
--

LOCK TABLES `schedule` WRITE;
/*!40000 ALTER TABLE `schedule` DISABLE KEYS */;
INSERT INTO `schedule` VALUES (1,'示例任务-关联场景-一次性','once','','no','scene','1','','1','示例场景','示例产品','not_started','','','','',NULL,NULL,'admin','2023-12-13 07:18:56',NULL,NULL),(2,'示例任务-关联数据-自定义执行','cron',' */2 * * * *','no','data','','示例-用户管理-新建用户','1','','示例产品','not_started','','','','',NULL,NULL,'admin','2023-12-13 07:19:36',NULL,NULL),(3,'示例任务-关联场景-每天0点和12点执行','day','','no','scene','1','','1','示例场景','示例产品','not_started','','0,12','','',NULL,NULL,'admin','2023-12-13 07:20:23',NULL,NULL),(4,'示例任务-关联场景-每周六和每周日0点和20点执行','week','','no','scene','1','','1','示例场景','示例产品','not_started','0,20','','6,7','',NULL,NULL,'admin','2023-12-13 07:21:16',NULL,NULL);
/*!40000 ALTER TABLE `schedule` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-01-05 20:55:26
