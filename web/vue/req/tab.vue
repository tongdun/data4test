<template>
  <div class="layout">
    <div class="layout-content">
      <Tabs :animated="false" :tab="tab.id" name="top" value="run" @on-click="setDomainMode">
        <Tab-pane :label="defLabel" name="def" tab="top">
          <i-col span="4">
            <Menu
                ref="side_menu"
                accordion
                theme="light"
                width="auto"
                @on-select="selectApiDetail"
                @on-open-change="selectDefModule"
            >
              <Submenu v-for="item in menuData" :key="item.name" ref="child" :name="item.name">
                <template slot="title">
                  <i :class="'iconfont '+item.icon"></i>
                  <span>{{item.title}}</span>
                </template>
                <template v-for="list1 in item.children">
                  <Submenu v-if="list1.children&&list1.children.length!==0" :name="list1.name">
                    <template slot="title">
                      <i :class="'iconfont '+'11'"></i>
                      <span>{{list1.title}}</span>
                    </template>
                    <MenuItem
                        v-for="list2 in list1.children"
                        :key="list2.name"
                        :name="list2.name">
                      {{list2.title}}
                    </MenuItem>
                  </Submenu>
                  <MenuItem
                      v-else
                      :name="list1.name"
                      class="noChildmenuitem"
                  >
                    <i :class="'iconfont '+'11'"></i>
                    &nbsp;&nbsp;&nbsp;
                    {{ list1.title }}
                  </MenuItem>
                </template>
              </Submenu>
            </Menu>
          </i-col>
          <i-col span="20">
            <Form>
              <Form-item>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Select v-model="apiDefSave.app" :placeholder="$t('api.appTips')" allow-create clearable filterable style="width:10%;" @on-create="onAddApp" @on-change="getAppData">
                      <Option v-for="item in appOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Input v-model="apiDefSave.module" :disabled="isSending" :placeholder="$t('api.moduleTips')" style="width:20%;"></Input>
                    <Input v-model="apiDefSave.apiDesc" :disabled="isSending" :placeholder="$t('api.defDescTips')" style="width:20%;"></Input>
                  </Col>
                </Row>
                <br>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Select v-model="apiDefSave.method" :placeholder="$t('api.methodTips')" allow-create clearable filterable style="width:9%;" @on-create="onAddMethod">
                      <Option v-for="item in defMethodOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Input v-model="apiDefSave.prefix" :placeholder="$t('api.prefixTips')" disabled style="width:10%;"></Input>
                    <Input v-model="apiDefSave.path" :disabled="isSending" :placeholder="$t('api.urlTips')" style="width:70%;"></Input>
                  </Col>
                  <Col :lg="3" :md="4" :sm="5" :xs="6" style="display:flex;justify-content:flex-end;">
                    <Button-group>
                      <Button :disabled="isSending" :icon="saveIcon" class="layout-button-left" type="primary" @click="onDefSave">{{$t('buttonTitle.apiSave')}}</Button>
                    </Button-group>
                  </Col>
                </Row>
              </Form-item>
              <Form-item>
                <Card style="width:100%;">
                  <div slot="title">{{$t('tabTitle.contentTab')}}</div>
                  <div>
                    <Tabs :animated="false" name="child-name" tab="def" type="line" value="Body">
                      <Tab-pane v-model="apiDefSave.pathVars" :label="defPathLabel" name="Path" tab="child-name">
                        <DefList :defListData="apiDefSave.pathVars" :isSending="isSending"></DefList>
                      </Tab-pane>
                      <Tab-pane v-model="apiDefSave.queryVars" :label="defQuerysLabel" name="Query" tab="child-name">
                        <DefList :defListData="apiDefSave.queryVars" :isSending="isSending"></DefList>
                      </Tab-pane>
                      <Tab-pane v-model="apiDefSave.bodyVars" :label="defBodyLabel" name="Body" tab="child-name">
                        <div style="padding-bottom:10px;">
                          <RadioGroup v-model="apiDefSave.bodyMode" size="large">
                            <Radio label="application/x-www-form-urlencoded">application/x-www-form-urlencoded</Radio>
                            <Radio label="application/json">application/json</Radio>
                            <Radio label="multipart/form-data">multipart/form-data</Radio>
                            <Radio label="raw">raw</Radio>
                          </RadioGroup>
                          <Select v-model="rawDefContentType" :disabled="!isDefBodyRaw" :transfer="true" style="width:220px;">
                            <Option value="x-www-form-urlencoded">Form-Data(x-www-form-urlencoded)</Option>
                            <Option value="json">JSON(application/json)</Option>
                            <Option value="form-data">Binary(multipart/form-data)</Option>
                            <Option value="text">Text</Option>
                            <Option value="text/plain">Text(text/plain)</Option>
                            <Option value="application/javascript">Javascript(application/javascript)</Option>
                            <Option value="application/xml">XML(application/xml)</Option>
                            <Option value="text/xml">XML(text/xml)</Option>
                            <Option value="text/html">HTML(text/html)</Option>
                          </Select>
                        </div>
                        <DefList v-show="!isDefBodyRaw" :defListData="apiDefSave.bodyVars" :isSending="isSending"></DefList>
                        <Input v-show="isDefBodyRaw" v-model="apiDefSave.bodyStr" :rows="10" type="textarea"></Input>
                      </Tab-pane>
                      <Tab-pane v-model="apiDefSave.headerVars" :label="defHeadersLabel" name="Header" tab="child-name">
                        <DefList :defListData="apiDefSave.headerVars" :isSending="isSending"></DefList>
                      </Tab-pane>
                      <Tab-pane v-model="apiDefSave.respVars" :label="defRespLabel" name="Resp" tab="child-name">
                        <DefList :defListData="apiDefSave.respVars" :isSending="isSending"></DefList>
                      </Tab-pane>
                    </Tabs>
                  </div>
                </Card>
              </Form-item>
            </Form>
          </i-col>
        </Tab-pane>

        <Tab-pane :label="runLabel" name="run" tab="top">
          <i-col span="4">
            <Menu
                ref="side_menu"
                accordion
                theme="light"
                width="auto"
                @on-select="selectRunApiDesc"
                @on-open-change="selectRunModule"
            >
              <Submenu v-for="item in menuData" :key="item.name" ref="child" :name="item.name">
                <template slot="title">
                  <i :class="'iconfont '+item.icon"></i>
                  <span>{{item.title}}</span>
                </template>
                <template v-for="list1 in item.children">
                  <Submenu v-if="list1.children&&list1.children.length!==0" :name="list1.name">
                    <template slot="title">
                      <i :class="'iconfont '+'11'"></i>
                      <span>{{list1.title}}</span>
                    </template>
                    <MenuItem
                        v-for="list2 in list1.children"
                        :key="list2.name"
                        :name="list2.name">
                      {{list2.title}}
                    </MenuItem>
                  </Submenu>
                  <MenuItem
                      v-else
                      :name="list1.name"
                      class="noChildmenuitem"
                  >
                    <i :class="'iconfont '+'11'"></i>
                    &nbsp;&nbsp;&nbsp;
                    {{ list1.title }}
                  </MenuItem>
                </template>
              </Submenu>
            </Menu>
          </i-col>
          <i-col span="20">
            <Form>
              <Form-item>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Select v-model="apiRunSave.app" :placeholder="$t('api.appTips')" allow-create clearable filterable style="width:10%;"  @on-create="onAddApp" @on-change="getAppData">
                      <Option v-for="item in appOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Input v-model="apiRunSave.module" :disabled="isSending" :placeholder="$t('api.moduleTips')" style="width:20%;"></Input>
                    <Input v-model="apiRunSave.apiDesc" :disabled="isSending" :placeholder="$t('api.defDescTips')" style="width:20%;"></Input>
                    <Select v-model="apiRunSave.dataDesc" :placeholder="$t('api.dataDescTips')" allow-create clearable filterable style="width:48%;" @on-create="onAddDataDesc" @on-change="getApiDataDetail">
                      <Option v-for="item in dataDescOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                  </Col>
                </Row>
                <br>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Select v-model="apiRunSave.method" :placeholder="$t('api.methodTips')" allow-create clearable filterable style="width:100px" @on-create="onAddMethod" @on-change="getMethodData">
                      <Option v-for="item in runMethodOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Select v-model="apiRunSave.prototype" :placeholder="$t('api.prototypeTips')" style="width:100px" value="http">
                      <Option v-for="item in prototypeOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Input v-model="apiRunSave.hostIp" :placeholder="$t('api.hostIpTips')" style="width:15%;"></Input>
                    <Input v-model="apiRunSave.prefix" :placeholder="$t('api.prefixTips')" style="width:8%;"></Input>
                    <Input v-model="apiRunSave.path" :placeholder="$t('api.urlTips')" style="width:30%;"></Input>
                    <Select v-model="apiRunSave.product" :placeholder="$t('api.envTips')" clearable filterable style="width:15%;" @on-change="getApiEnv">-->
                      <Option v-for="item in depDataTableData" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                  </Col>

                  <Col :lg="3" :md="4" :sm="5" :xs="6" style="display:flex;justify-content:flex-end;">
                    <Button-group>
                      <Button :loading="isSending" icon="android-send" type="primary" @click="apiRun">{{$t('buttonTitle.dataSend')}}</Button>&nbsp;
                    </Button-group>
                    <Button-group>
                      <Button :disabled="isSending" :icon="saveIcon" type="primary" @click="onApiRunSave">{{$t('buttonTitle.dataSave')}}</Button>
                    </Button-group>
                  </Col>

                </Row>
              </Form-item>
              <Form-item>
                <Card style="width:100%;">
                  <div slot="title">{{$t('tabTitle.contentTab')}}</div>
                  <div>
                    <Tabs :animated="false" name="subRun" tab="run" type="line" value="Body">
                      <Tab-pane :label="pathLabel" name="Path" tab="subRun">
                        <RunList :isSending="isSending" :runListData="apiRunSave.pathVars" ></RunList>
                      </Tab-pane>
                      <Tab-pane :label="querysLabel" name="Query" tab="subRun">
                        <RunList :isSending="isSending" :runListData="apiRunSave.queryVars"></RunList>
                      </Tab-pane>
                      <Tab-pane :label="bodyLabel" name="Body" tab="subRun">
                        <div style="padding-bottom:10px;">
                          <RadioGroup v-model="apiRunSave.bodyMode" size="large">
                            <Radio label="application/x-www-form-urlencoded">application/x-www-form-urlencoded</Radio>
                            <Radio label="application/json">application/json</Radio>
                            <Radio label="multipart/form-data">multipart/form-data</Radio>
                            <Radio label="raw">raw</Radio>
                          </RadioGroup>
                          <Select v-model="rawRunContentType" :disabled="!isRunBodyRaw" :transfer="true" style="width:220px;">
                            <Option value="x-www-form-urlencoded">Form-Data(x-www-form-urlencoded)</Option>
                            <Option value="json">JSON(application/json)</Option>
                            <Option value="form-data">Binary(multipart/form-data)</Option>
                            <Option value="text">Text</Option>
                            <Option value="text/plain">Text(text/plain)</Option>
                            <Option value="application/javascript">Javascript(application/javascript)</Option>
                            <Option value="application/xml">XML(application/xml)</Option>
                            <Option value="text/xml">XML(text/xml)</Option>
                            <Option value="text/html">HTML(text/html)</Option>
                          </Select>
                        </div>
                        <RunList v-show="!isRunBodyRaw" :isSending="isSending" :runListData="apiRunSave.bodyVars"></RunList>
                        <Input v-show="isRunBodyRaw" v-model="apiRunSave.bodyStr" :rows="10" type="textarea"></Input>
                      </Tab-pane>
                      <Tab-pane :label="headersLabel" name="Headers" tab="subRun">
                        <RunList :isSending="isSending" :runListData="apiRunSave.headerVars"></RunList>
                      </Tab-pane>
                      <Tab-pane :label="respsLabel" name="Resps" tab="subRun">
                        <RunList :isSending="isSending" :runListData="apiRunSave.respVars"></RunList>
                      </Tab-pane>
                      <Tab-pane :label="actionsLabel" name="Actions" tab="subRun">
                        <ActionList :actionListData="apiRunSave.actions" :isSending="isSending"></ActionList>
                      </Tab-pane>
                      <Tab-pane :label="assertsLabel" name="Asserts" tab="subRun">
                        <AssertList :assertListData="apiRunSave.asserts" :isSending="isSending"></AssertList>
                      </Tab-pane>
                      <Tab-pane :label="preApisLabel" name="PreApis" tab="subRun">
                        <RelatedApiList :allDataFile="allDataFile" :isSending="isSending" :relatedApiListData="apiRunSave.preApis"></RelatedApiList>
                      </Tab-pane>
                      <Tab-pane :label="postApisLabel" name="PostApis" tab="subRun">
                        <RelatedApiList :allDataFile="allDataFile" :isSending="isSending" :relatedApiListData="apiRunSave.postApis"></RelatedApiList>
                      </Tab-pane>
                      <Tab-pane :label="otherConfigLabel" name="OtherConfig" tab="subRun">
                        <OtherConfigList :isSending="isSending" :otherConfigListData="apiRunSave.otherConfigs"></OtherConfigList>
                      </Tab-pane>
                    </Tabs>
                  </div>
                </Card>
              </Form-item>
              <Form-item v-show="isResponse">
                <RespInfo :requestData="reqDataRespList"></RespInfo>
                <RequestInfo :requestData="reqDataRespList"></RequestInfo>
              </Form-item>
            </Form>
          </i-col>
        </Tab-pane>

        <Tab-pane :label="dataLabel" name="data" tab="top">
          <i-col span="4">
            <Menu
                ref="side_menu"
                accordion
                theme="light"
                width="auto"
                @on-select="selectData"
                @on-open-change="selectAppData"
            >
              <Submenu v-for="item in dataMenu" :key="item.name" ref="child" :name="item.name">
                <template slot="title">
                  <span>{{item.title}}</span>
                </template>
                <MenuItem
                    v-for="subItem in item.children"
                    :key="subItem"
                    :name="subItem">
                  {{subItem}}
                </MenuItem>
              </Submenu>
            </Menu>
          </i-col>
          <i-col span="20">
            <Form>
              <Form-item>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Select v-model="dataRunSave.app" :placeholder="$t('api.appTips')" allow-create clearable filterable style="width:10%;" @on-create="onAddApp" @on-change="getDataApp">
                      <Option v-for="item in appOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Input v-model="dataRunSave.module" :disabled="isSending" :placeholder="$t('api.moduleTips')" style="width:20%;"></Input>
                    <Input v-model="dataRunSave.apiDesc" :disabled="isSending" :placeholder="$t('api.defDescTips')" style="width:20%;"></Input>
                    <Select v-model="dataRunSave.dataDesc" :placeholder="$t('api.dataDescTips')" allow-create clearable filterable style="width:48%;" @on-create="onAddDataDesc" @on-change="getDataByName">
                      <Option v-for="item in dataDescOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                  </Col>
                </Row>
                <br>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Select v-model="dataRunSave.method" :placeholder="$t('api.methodTips')" allow-create clearable filterable style="width:100px" @on-create="onAddMethod">
                      <Option v-for="item in runMethodOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Select v-model="dataRunSave.prototype" :placeholder="$t('api.prototypeTips')" style="width:100px" value="http">
                      <Option v-for="item in prototypeOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Input v-model="dataRunSave.hostIp" :placeholder="$t('api.hostIpTips')" style="width:15%;"></Input>
                    <Input v-model="dataRunSave.prefix" :placeholder="$t('api.prefixTips')" style="width:8%;"></Input>
                    <Input v-model="dataRunSave.path" :placeholder="$t('api.urlTips')" style="width:30%;"></Input>

                    <Select v-model="dataRunSave.product" :placeholder="$t('api.envTips')" clearable  filterable style="width:15%;" @on-change="getDataEnv">-->
                      <Option v-for="item in depDataTableData" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                  </Col>

                  <Col :lg="3" :md="4" :sm="5" :xs="6" style="display:flex;justify-content:flex-end;">
                    <Button-group>
                      <Button :loading="isSending" icon="android-send" type="primary" @click="dataRun">{{$t('buttonTitle.dataSend')}}</Button>&nbsp;
                    </Button-group>
                    <Button-group>
                      <Button :disabled="isSending" :icon="saveIcon" type="primary" @click="onDataSave">{{$t('buttonTitle.dataSave')}}</Button>
                    </Button-group>
                  </Col>

                </Row>
              </Form-item>
              <Form-item>
                <Card style="width:100%;">
                  <div slot="title">{{$t('tabTitle.contentTab')}}</div>
                  <div>
                    <Tabs :animated="false" name="subRun" tab="run" type="line" value="Body">
                      <Tab-pane :label="dataPathLabel" name="Path" tab="subRun">
                        <DataRunList :isSending="isSending" :runListData="dataRunSave.pathVars" ></DataRunList>
                      </Tab-pane>
                      <Tab-pane :label="dataQuerysLabel" name="Query" tab="subRun">
                        <DataRunList :isSending="isSending" :runListData="dataRunSave.queryVars"></DataRunList>
                      </Tab-pane>
                      <Tab-pane :label="dataBodyLabel" name="Body" tab="subRun">
                        <div style="padding-bottom:10px;">
                          <RadioGroup v-model="dataRunSave.bodyMode" size="large">
                            <Radio label="application/x-www-form-urlencoded">application/x-www-form-urlencoded</Radio>
                            <Radio label="application/json">application/json</Radio>
                            <Radio label="multipart/form-data">multipart/form-data</Radio>
                            <Radio label="raw">raw</Radio>
                          </RadioGroup>
                          <Select v-model="dataRunSave.bodyMode" :disabled="!isDataBodyRaw" :transfer="true" style="width:220px;">
                            <Option value="x-www-form-urlencoded">Form-Data(x-www-form-urlencoded)</Option>
                            <Option value="json">JSON(application/json)</Option>
                            <Option value="form-data">Binary(multipart/form-data)</Option>
                            <Option value="text">Text</Option>
                            <Option value="text/plain">Text(text/plain)</Option>
                            <Option value="application/javascript">Javascript(application/javascript)</Option>
                            <Option value="application/xml">XML(application/xml)</Option>
                            <Option value="text/xml">XML(text/xml)</Option>
                            <Option value="text/html">HTML(text/html)</Option>
                          </Select>
                        </div>
                        <DataRunList v-show="!isDataBodyRaw" :isSending="isSending" :runListData="dataRunSave.bodyVars"></DataRunList>
                        <Input v-show="isDataBodyRaw" v-model="dataRunSave.bodyStr" :rows="10" type="textarea"></Input>
                      </Tab-pane>
                      <Tab-pane :label="dataHeadersLabel" name="Headers" tab="subRun">
                        <DataRunList :isSending="isSending" :runListData="dataRunSave.headerVars"></DataRunList>
                      </Tab-pane>
                      <Tab-pane :label="dataRespLabel" name="Resps" tab="subRun">
                        <DataRunList :isSending="isSending" :runListData="dataRunSave.respVars"></DataRunList>
                      </Tab-pane>
                      <Tab-pane :label="dataActionsLabel" name="Actions" tab="subRun">
                        <ActionList :actionListData="dataRunSave.actions" :isSending="isSending"></ActionList>
                      </Tab-pane>
                      <Tab-pane :label="dataAssertsLabel" name="Asserts" tab="subRun">
                        <AssertList :assertListData="dataRunSave.asserts" :isSending="isSending"></AssertList>
                      </Tab-pane>
                      <Tab-pane :label="dataPreApisLabel" name="PreApis" tab="subRun">
                        <RelatedApiList :allDataFile="allDataFile" :isSending="isSending" :relatedApiListData="dataRunSave.preApis"></RelatedApiList>
                      </Tab-pane>
                      <Tab-pane :label="dataPostApisLabel" name="PostApis" tab="subRun">
                        <RelatedApiList :allDataFile="allDataFile" :isSending="isSending" :relatedApiListData="dataRunSave.postApis"></RelatedApiList>
                      </Tab-pane>
                      <Tab-pane :label="dataOtherConfigLabel" name="OtherConfig" tab="subRun">
                        <OtherConfigList :isSending="isSending" :otherConfigListData="dataRunSave.otherConfigs"></OtherConfigList>
                      </Tab-pane>
                    </Tabs>
                  </div>
                </Card>
              </Form-item>
              <Form-item v-show="isResponse">
                <RequestInfo :requestData="dataModeReqDataRespList"></RequestInfo>
                <RespInfo :requestData="dataModeReqDataRespList"></RespInfo>
              </Form-item>
            </Form>
          </i-col>
        </Tab-pane>

        <Tab-pane :label="sceneLabel" name="scene" tab="top">
          <i-col span="4">
            <Menu
                ref="side_menu"
                accordion
                theme="light"
                width="auto"
                @on-select="selectScene"
                @on-open-change="selectProductScene"
            >
              <Submenu v-for="item in sceneMenu" :key="item.name" ref="child" :name="item.name">
                <template slot="title">
                  <span>{{item.title}}</span>
                </template>
                <MenuItem
                    v-for="subItem in item.children"
                    :key="subItem"
                    :name="subItem">
                  {{subItem}}
                </MenuItem>
              </Submenu>
            </Menu>
          </i-col>
          <i-col span="20">
            <Form>
              <Form-item>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Select v-model="sceneSave.name" :placeholder="$t('scene.nameTips')" allow-create clearable filterable style="width:55%;" @on-create="onAddPlaybook"  @on-change="getSceneByName">
                      <Option v-for="item in playbookOptions" :key="item" :value="item" allow-create>{{ item }}</Option>
                    </Select>
                    <Select v-model="sceneSave.type" :placeholder="$t('scene.typeTips')" clearable filterable style="width:10%;">
                      <Option v-for="item in sceneTypeOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Input v-model="sceneSave.runNum" :disabled="isSending" :placeholder="$t('scene.runNumTips')" clearable  style="width:10%;" type='number'></Input>
                    <Select v-model="sceneSave.product" :placeholder="$t('api.envTips')" clearable  filterable style="width:20%;" @on-change="getSceneEnv">
                      <Option v-for="item in depDataTableData" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                  </Col>
                  <Col :lg="3" :md="4" :sm="5" :xs="6" style="display:flex;justify-content:flex-end;">
                    <Button-group>
                      <Button :loading="isSending" icon="android-send" type="primary" @click="sceneRun">{{$t('buttonTitle.dataSend')}}</Button>&nbsp;
                    </Button-group>
                    <Button-group>
                      <Button :disabled="isSending" :icon="saveIcon" type="primary" @click="onSceneSave">{{$t('buttonTitle.sceneSave')}}</Button>
                    </Button-group>
                  </Col>

                </Row>

              </Form-item>
              <Form-item>
                <Card style="width:100%;">
                  <div slot="title">{{$t('tabTitle.contentTab')}}</div>
                  <div>
                    <Tabs :animated="false" name="subScene" tab="scene" type="line" value="dataFile">
                      <Tab-pane :label="datasLabel" name="dataFile" tab="subScene">
                        <RelatedApiList :allDataFile="allDataFile" :isSending="isSending" :relatedApiListData="sceneSave.dataList"></RelatedApiList>
                      </Tab-pane>
                    </Tabs>
                  </div>
                </Card>
              </Form-item>
              <Form-item v-show="isResponse">
                <RequestSceneInfo :requestScene="sceneModeReqDataRespList"></RequestSceneInfo>
              </Form-item>
            </Form>
          </i-col>
        </Tab-pane>

        <Tab-pane :label="historyLabel" name="history" tab="top">
          <i-col span="4">
            <Menu
                ref="side_menu"
                accordion
                theme="light"
                width="auto"
                @on-select="selectHistory"
                @on-open-change="selectHistoryDate"
            >
              <Submenu v-for="item in historyMenu" :key="item.name" ref="child" :name="item.name">
                <template slot="title">
                  <span>{{item.title}}</span>
                </template>
                <MenuItem
                    v-for="subItem in item.children"
                    :key="subItem"
                    :name="subItem">
                  {{subItem}}
                </MenuItem>
              </Submenu>
            </Menu>
          </i-col>
          <i-col span="20">
            <Form>
              <Form-item>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Input v-model="historyRunSave.app" :disabled="isSending" :placeholder="$t('api.appTips')" style="width:10%;"></Input>
                    <Input v-model="historyRunSave.module" :disabled="isSending" :placeholder="$t('api.moduleTips')" style="width:20%;"></Input>
                    <Input v-model="historyRunSave.apiDesc" :disabled="isSending" :placeholder="$t('api.defDescTips')" style="width:20%;"></Input>
                    <Input v-model="historyRunSave.dataDesc" :disabled="isSending" :placeholder="$t('api.dataDescTips')" style="width:48%;"></Input>
                  </Col>
                </Row>
                <br>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Select v-model="historyRunSave.method" :placeholder="$t('api.methodTips')" allow-create clearable filterable style="width:100px" @on-create="onAddMethod">
                      <Option v-for="item in runMethodOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Select v-model="historyRunSave.prototype" :placeholder="$t('api.prototypeTips')" style="width:100px" value="http">
                      <Option v-for="item in prototypeOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Input v-model="historyRunSave.host" :placeholder="$t('api.hostIpTips')" style="width:15%;"></Input>
                    <Input v-model="historyRunSave.prefix" :placeholder="$t('api.prefixTips')" style="width:8%;"></Input>
                    <Input v-model="historyRunSave.path" :placeholder="$t('api.urlTips')" style="width:30%;"></Input>
                    <Select v-model="historyRunSave.product" :placeholder="$t('api.envTips')" clearable  filterable style="width:15%;" @on-change="getHistoryEnv">-->
                      <Option v-for="item in depDataTableData" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                  </Col>

                  <Col :lg="3" :md="4" :sm="5" :xs="6" style="display:flex;justify-content:flex-end;">
                    <Button-group>
                      <Button :loading="isSending" icon="android-send" type="primary" @click="historyRun">{{$t('buttonTitle.dataSend')}}</Button>&nbsp;
                    </Button-group>

                    <Button-group>
                      <Button :disabled="isSending" :icon="saveIcon" type="primary" @click="onHistorySave">{{$t('buttonTitle.dataSave')}}</Button>
                    </Button-group>
                  </Col>

                </Row>
              </Form-item>
              <Form-item>
                <Card style="width:100%;">
                  <div slot="title">{{$t('tabTitle.contentTab')}}</div>
                  <div>
                    <Tabs :animated="false" name="subRun" tab="run" type="line" value="Body">
                      <Tab-pane :label="historyPathLabel" name="Path" tab="subRun">
                        <DataRunList :isSending="isSending" :runListData="historyRunSave.pathVars" ></DataRunList>
                      </Tab-pane>
                      <Tab-pane :label="historyQuerysLabel" name="Query" tab="subRun">
                        <DataRunList :isSending="isSending" :runListData="historyRunSave.queryVars"></DataRunList>
                      </Tab-pane>
                      <Tab-pane :label="historyBodyLabel" name="Body" tab="subRun">
                        <div style="padding-bottom:10px;">
                          <RadioGroup v-model="historyRunSave.bodyMode" size="large">
                            <Radio label="application/x-www-form-urlencoded">application/x-www-form-urlencoded</Radio>
                            <Radio label="application/json">application/json</Radio>
                            <Radio label="multipart/form-data">multipart/form-data</Radio>
                            <Radio label="raw">raw</Radio>
                          </RadioGroup>
                          <Select v-model="rawHistoryContentType" :disabled="!isHistoryBodyRaw" :transfer="true" style="width:220px;">
                            <Option value="x-www-form-urlencoded">Form-Data(x-www-form-urlencoded)</Option>
                            <Option value="json">JSON(application/json)</Option>
                            <Option value="form-data">Binary(multipart/form-data)</Option>
                            <Option value="text">Text</Option>
                            <Option value="text/plain">Text(text/plain)</Option>
                            <Option value="application/javascript">Javascript(application/javascript)</Option>
                            <Option value="application/xml">XML(application/xml)</Option>
                            <Option value="text/xml">XML(text/xml)</Option>
                            <Option value="text/html">HTML(text/html)</Option>
                          </Select>
                        </div>
                        <DataRunList v-show="!isHistoryBodyRaw" :isSending="isSending" :runListData="historyRunSave.bodyVars"></DataRunList>
                        <Input v-show="isHistoryBodyRaw" v-model="historyRunSave.bodyStr" :rows="10" type="textarea"></Input>
                      </Tab-pane>
                      <Tab-pane :label="historyHeadersLabel" name="Headers" tab="subRun">
                        <DataRunList :isSending="isSending" :runListData="historyRunSave.headerVars"></DataRunList>
                      </Tab-pane>
                      <Tab-pane :label="historyRespsLabel" name="Resps" tab="subRun">
                        <DataRunList :isSending="isSending" :runListData="historyRunSave.respVars"></DataRunList>
                      </Tab-pane>
                      <Tab-pane :label="historyActionsLabel" name="Actions" tab="subRun">
                        <ActionList :actionListData="historyRunSave.actions" :isSending="isSending"></ActionList>
                      </Tab-pane>
                      <Tab-pane :label="historyAssertsLabel" name="Asserts" tab="subRun">
                        <AssertList :assertListData="historyRunSave.asserts" :isSending="isSending"></AssertList>
                      </Tab-pane>
                      <Tab-pane :label="historyPreApisLabel" name="PreApis" tab="subRun">
                        <RelatedApiList :allDataFile="allDataFile" :isSending="isSending" :relatedApiListData="historyRunSave.preApis"></RelatedApiList>
                      </Tab-pane>
                      <Tab-pane :label="historyPostApisLabel" name="PostApis" tab="subRun">
                        <RelatedApiList :allDataFile="allDataFile" :isSending="isSending" :relatedApiListData="historyRunSave.postApis"></RelatedApiList>
                      </Tab-pane>
                      <Tab-pane :label="historyOtherConfigLabel" name="OtherConfig" tab="subRun">
                        <OtherConfigList :isSending="isSending" :otherConfigListData="historyRunSave.otherConfigs"></OtherConfigList>
                      </Tab-pane>
                    </Tabs>
                  </div>
                </Card>
              </Form-item>
              <Form-item v-show="isHistoryResponse">
                <RequestInfo :requestData="historyModeReqDataRespList"></RequestInfo>
                <RespInfo :requestData="historyModeReqDataRespList"></RespInfo>
              </Form-item>
            </Form>
          </i-col>
        </Tab-pane>

        <Tab-pane :label="sceneHistoryLabel" name="sceneHistory" tab="top">
          <i-col span="4">
            <Menu
                ref="side_menu"
                accordion
                theme="light"
                width="auto"
                @on-select="selectSceneHistory"
                @on-open-change="selectSceneHistoryDate"
            >
              <Submenu v-for="item in sceneHistoryMenu" :key="item.name" ref="child" :name="item.name">
                <template slot="title">
                  <span>{{item.title}}</span>
                </template>
                <MenuItem
                    v-for="subItem in item.children"
                    :key="subItem"
                    :name="subItem">
                  {{subItem}}
                </MenuItem>
              </Submenu>
            </Menu>
          </i-col>
          <i-col span="20">
            <Form>
              <Form-item>
                <Row type="flex">
                  <Col :lg="21" :md="20" :sm="19" :xs="18">
                    <Input v-model="sceneHistorySave.name" :disabled="isSending" :placeholder="$t('scene.nameTips')" style="width:55%;"></Input>
                    <Select v-model="sceneHistorySave.type" :placeholder="$t('scene.typeTips')" clearable filterable style="width:10%;">
                      <Option v-for="item in sceneTypeOptions" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                    <Input v-model="sceneHistorySave.runNum" :disabled="isSending" :placeholder="$t('scene.runNumTips')" clearable  style="width:10%;" type='number'></Input>
                    <Select v-model="sceneHistorySave.product" :placeholder="$t('api.envTips')" clearable  filterable style="width:20%;" @on-change="getSceneEnv">
                      <Option v-for="item in depDataTableData" :key="item" :value="item">{{ item }}</Option>
                    </Select>
                  </Col>
                  <Col :lg="3" :md="4" :sm="5" :xs="6" style="display:flex;justify-content:flex-end;">
                    <Button-group>
                      <Button :loading="isSending" icon="android-send" type="primary" @click="sceneHistoryRun">{{$t('buttonTitle.dataSend')}}</Button>&nbsp;
                    </Button-group>
                    <Button-group>
                      <Button :disabled="isSending" :icon="saveIcon" type="primary" @click="onSceneHistorySave">{{$t('buttonTitle.sceneSave')}}</Button>
                    </Button-group>
                  </Col>

                </Row>

              </Form-item>
              <Form-item>
                <Card style="width:100%;">
                  <div slot="title">{{$t('tabTitle.contentTab')}}</div>
                  <div>
                    <Tabs :animated="false" name="subScene" tab="scene" type="line" value="dataFile">
                      <Tab-pane :label="datasLabel" name="dataFile" tab="subScene">
                        <historyApiList :allDataFile="allDataFile" :historyApiListData="sceneHistorySave.dataList" :isSending="isSending"></historyApiList>
                      </Tab-pane>
                    </Tabs>
                  </div>
                </Card>
              </Form-item>
              <Form-item v-show="isSceneHistoryResponse">
                <RequestSceneInfo :requestScene="sceneHistoryRespList"></RequestSceneInfo>
              </Form-item>
            </Form>
          </i-col>
        </Tab-pane>

      </Tabs>
    </div>
  </div>
</template>
<style scoped>
.layout{
  border: 1px solid #d7dde4;
  background: #f5f7f9;
}

.layout-content{
  min-height: 200px;
  margin: 15px;
  overflow: hidden;
  background: #fff;
  border-radius: 4px;
}


</style>
<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import _, {forEach} from 'lodash'
import moment from 'moment'
import API from '../../ts/api'
import ReqPre from './pre.vue'
import ReqReport from './report.vue'
import { CreateElement } from 'vue/types/vue'
import DefList from './defList.vue'
import RunList from "./list.vue"
import DataRunList from "./dataReq.vue"
import ActionList from './action.vue'
import AssertList from './assert.vue'
import EnvList from './env.vue'
import OtherConfigList from './otherConfig.vue'
import debounce from 'lodash/debounce'
import RelatedApiList from './depApi.vue'
import RespInfo from './response.vue'
import RequestInfo from './request.vue'
import RequestSceneInfo from './sceneResponse.vue'
import AdvancedInfo from './advance.vue'
import Test from './test.vue'
import Router from "../../ts/router";
import App from "../app.vue";
import historyApiList from "./historyPlaybook.vue";

@Component({
  components: {
    ReqPre,
    ReqReport,
    DefList,
    RunList,
    DataRunList,
    ActionList,
    AssertList,
    EnvList,
    OtherConfigList,
    RelatedApiList,
    RespInfo,
    RequestInfo,
    RequestSceneInfo,
    AdvancedInfo,
    historyApiList
  }
})
export default class Tab extends Vue {
  @Prop() tab: Req.TabModel
  menuData: any[] = []
  dataMenu: any[] = []
  sceneMenu: any[] = []
  historyMenu: any[] = []
  sceneHistoryMenu: any[] = []
  appOptions: string[] = []
  playbookOptions: string[] = []

  defMethodOptions: string[] = ["get", "post", "put", "delete"]
  defApiOptions: string[] = []
  defModuleOptions: string[] = []
  defApiDescOptions: string[] = []

  runMethodOptions: string[] = ["get", "post", "put", "delete"]
  prototypeOptions: string[] = ["http", "https"]
  runApiOptions: string[] = []
  runModuleOptions: string[] = []
  runApiDescOptions: string[] = []
  dataDescOptions: string[] = []
  depDataTableData: Req.DepDataListModel[] = []  //可以从其他页面捞来数据

  sceneTypeOptions: string[] = ["默认", "比较"]
  allDataFile: string[] = []

  domainMode = "run"
  advancedIcon = 'ios-arrow-down'
  isAdvanced = true
  isResponse = false
  isHistoryResponse = true
  isResquest = false

  isSceneHistoryResponse = true

  saveIcon = 'save'
  isSaved = true
  previewIcon = 'ios-arrow-down'
  previewType = 'primary'
  isPreview = true
  isSending = false
  reqMode: 'Path' | 'Query' | 'Body' | 'Header' | '断言' | '前置接口' | '后置接口' | '环境'
  contentType = ''
  params: Req.RunListModel[] = []

  bodyMode: 'application/x-www-form-urlencoded' | 'application/json' | 'multipart/form-data' | 'raw' = 'application/x-www-form-urlencoded'

  actions:  Req.RunListModel[] = []
  asserts: Req.RunListModel[] = []
  preApis: Req.RunListModel[] = []
  postApis: Req.RunListModel[] = []
  querys:  Req.RunListModel[] = []
  bodys: Req.RunListModel[] = []
  body: string = ''

  selectDate: string = ""
  sceneSelectDate: string = ""

  varDef: Req.DefListModel = {
    name: '',
    valueType: 'string',
    isMust: 'no',
    egValue: '',
    desc: ''
  }

  reqN: number = 10
  reqC: number = 10
  reqTimeout: number = 20

  req: Req.RequestModel = {
    n: 0,
    c: 0,
    timeout: 0,
    method: 'GET',
    url: '',
    headers: {},
    body: ''
  }

  sceneSave: Req.SceneModel = {
    product: '',
    name: '',
    dataList: [],
    type: '默认',
    runNum: 1
  }

  sceneHistorySave: Req.SceneHistoryModel = {
    product: '',
    name: '',
    dataList: [],
    type: '默认',
    runNum: 1,
    result: "",
    failReason: "",
    lastFile: ""
  }

  singFile: Req.RelatedApiListModel = {
    dataFile: ''
  }
  apiDefSave: Req.ApiDefSaveModel = {
    app: '',
    module: '',
    apiDesc: '',
    prefix: '',
    method: '',
    path: '',
    bodyMode: 'application/x-www-form-urlencoded',
    bodyStr: '',
    pathVars: [],
    queryVars: [],
    bodyVars: [],
    headerVars:[],
    respVars: [],
  }

  apiRunSave: Req.ApiRunSaveModel = {
    app: '',
    module: '',
    apiDesc: '',
    dataDesc: '',
    prefix: '',
    method: '',
    prototype: '',
    hostIp: '',
    path: '',
    product: '',
    bodyMode: 'application/x-www-form-urlencoded',
    bodyStr: '',
    pathVars: [],
    queryVars: [],
    bodyVars: [],
    headerVars:[],
    respVars: [],
    actions: [],
    asserts: [],
    preApis: [],
    postApis: [],
    otherConfigs: []
  }

  historyRunSave: Req.ApiHistorySaveModel = {
    app: '',
    module: '',
    apiDesc: '',
    dataDesc: '',
    prefix: '',
    method: '',
    prototype: '',
    host: '',
    path: '',
    product: '',
    fileName: '',
    bodyMode: 'application/x-www-form-urlencoded',
    bodyStr: '',
    pathVars: [],
    queryVars: [],
    bodyVars: [],
    headerVars:[],
    respVars: [],
    actions: [],
    asserts: [],
    preApis: [],
    postApis: [],
    otherConfigs: [],
    output: ''
    // response: '',
    // request: '',
    // url: '',
    // header: '',
    // testResult: '',
    // failReason: ''
  }

  dataRunSave: Req.ApiRunSaveModel = {
    app: '',
    module: '',
    apiDesc: '',
    dataDesc: '',
    prefix: '',
    method: '',
    prototype: '',
    hostIp: '',
    path: '',
    product: '',
    bodyMode: 'application/x-www-form-urlencoded',
    pathVars: [],
    queryVars: [],
    bodyVars: [],
    bodyStr: '',
    headerVars:[],
    respVars: [],
    actions: [],
    asserts: [],
    preApis: [],
    postApis: [],
    otherConfigs: []
  }

  start: number = 0

  reqDataRespList: Req.ReqDataRespModel = {
    response: "",
    url: "",
    header: "",
    request: "",
    testResult: "",
    failReason: "",
    output: ""
  }

  dataModeReqDataRespList: Req.ReqDataRespModel = {
    response: "",
    url: "",
    header: "",
    request: "",
    testResult: "",
    failReason: "",
    output: "",
  }

  AdvancedModelList: Req.AdvancedModel = {
    version: "",
    apiId: "",
    isParallel: "",
    isUseEnvConfig: ""
  }

  historyModeReqDataRespList: Req.ReqDataRespModel = {
    output: "",
    response: "",
    url: "",
    header: "",
    request: "",
    testResult: "",
    failReason: ""
  }

  sceneModeReqDataRespList: Req.ReqSceneRespModel = {
    lastDataFile: "",
    testResult: "",
    failReason: ""
  }

  sceneHistoryRespList: Req.ReqSceneRespModel = {
    lastDataFile: "",
    testResult: "",
    failReason: ""
  }

  rawDefContentType: string = 'json'
  rawRunContentType: string = 'json'
  rawDataContentType: string = 'json'
  rawHistoryContentType: string = 'json'

  showParams: boolean = true

  get querysLabel(): string {
    let label = 'Query'

    if (this.apiRunSave.queryVars.length > 0) {
      label += '(' + this.apiRunSave.queryVars.length.toString() + ')'
    }

    return label
  }

  get defQuerysLabel(): string {
    let label = 'Query'

    if (this.apiDefSave.queryVars.length > 0) {
      label += '(' + this.apiDefSave.queryVars.length.toString() + ')'
    }

    return label
  }

  get dataQuerysLabel(): string {
    let label = 'Query'

    if (this.dataRunSave.queryVars.length > 0) {
      label += '(' + this.dataRunSave.queryVars.length.toString() + ')'
    }

    return label
  }

  get historyQuerysLabel(): string {
    let label = 'Query'

    if (this.historyRunSave.queryVars.length > 0) {
      label += '(' + this.historyRunSave.queryVars.length.toString() + ')'
    }

    return label
  }

  get pathLabel(): string {
    let label = 'Path'

    if (this.apiRunSave.pathVars.length > 0) {
      label += '(' + this.apiRunSave.pathVars.length.toString() + ')'
    }

    return label
  }

  get defPathLabel(): string {
    let label = 'Path'

    if (this.apiDefSave.pathVars.length > 0) {
      label += '(' + this.apiDefSave.pathVars.length.toString() + ')'
    }

    return label
  }

  get dataPathLabel(): string {
    let label = 'Path'

    if (this.dataRunSave.pathVars.length > 0) {
      label += '(' + this.dataRunSave.pathVars.length.toString() + ')'
    }

    return label
  }

  get historyPathLabel(): string {
    let label = 'Path'

    if (this.historyRunSave.pathVars.length > 0) {
      label += '(' + this.historyRunSave.pathVars.length.toString() + ')'
    }

    return label
  }

  get defLabel(): string {
    let label = '定义域'
    return label
  }

  get runLabel(): string {
    let label = '运行域'
    return label
  }

  get dataLabel(): string {
    let label = '数据域'
    return label
  }

  get sceneLabel(): string {
    let label = '场景域'
    return label
  }

  get historyLabel(): string {
    let label = '数据历史域'
    return label
  }

  get sceneHistoryLabel(): string {
    let label = '场景历史域'
    return label
  }

  get bodyLabel(): string {
    let label = 'Body'

    if (this.apiRunSave.bodyVars.length > 0) {
      label += '(' + this.apiRunSave.bodyVars.length.toString() + ')'
    }

    return label
  }

  get defBodyLabel(): string {
    let label = 'Body'

    if (this.apiDefSave.bodyVars.length > 0) {
      label += '(' + this.apiDefSave.bodyVars.length.toString() + ')'
    }

    return label
  }

  get dataBodyLabel(): string {
    let label = 'Body'

    if (this.dataRunSave.bodyVars.length > 0) {
      label += '(' + this.dataRunSave.bodyVars.length.toString() + ')'
    }

    return label
  }

  get historyBodyLabel(): string {
    let label = 'Body'

    if (this.historyRunSave.bodyVars.length > 0) {
      label += '(' + this.historyRunSave.bodyVars.length.toString() + ')'
    }

    return label
  }

  get headersLabel(): string {
    let label = 'Header'

    if (this.apiRunSave.headerVars.length > 0) {
      label += '(' + this.apiRunSave.headerVars.length.toString() + ')'
    }

    return label
  }

  get defHeadersLabel(): string {
    let label = 'Header'

    if (this.apiDefSave.headerVars.length > 0) {
      label += '(' + this.apiDefSave.headerVars.length.toString() + ')'
    }

    return label
  }

  get dataHeadersLabel(): string {
    let label = 'Header'

    if (this.dataRunSave.headerVars.length > 0) {
      label += '(' + this.dataRunSave.headerVars.length.toString() + ')'
    }

    return label
  }

  get historyHeadersLabel(): string {
    let label = 'Header'

    if (this.historyRunSave.headerVars.length > 0) {
      label += '(' + this.historyRunSave.headerVars.length.toString() + ')'
    }

    return label
  }

  get respLabel(): string {
    let label = 'Resp'

    if (this.apiRunSave.respVars.length > 0) {
      label += '(' + this.apiRunSave.respVars.length.toString() + ')'
    }

    return label
  }

  get defRespLabel(): string {
    let label = 'Resp'

    if (this.apiDefSave.respVars.length > 0) {
      label += '(' + this.apiDefSave.respVars.length.toString() + ')'
    }

    return label
  }

  get dataRespLabel(): string {
    let label = 'Resp'

    if (this.dataRunSave.respVars.length > 0) {
      label += '(' + this.dataRunSave.respVars.length.toString() + ')'
    }

    return label
  }

  get historyRespLabel(): string {
    let label = 'Resp'

    if (this.historyRunSave.respVars.length > 0) {
      label += '(' + this.historyRunSave.respVars.length.toString() + ')'
    }

    return label
  }

  get actionsLabel(): string {
    let label = '动作'

    if (this.apiRunSave.actions.length > 0) {
      label += '(' + this.apiRunSave.actions.length.toString() + ')'
    }

    return label
  }

  get assertsLabel(): string {
    let label = '断言'

    if (this.apiRunSave.asserts.length > 0) {
      label += '(' + this.apiRunSave.asserts.length.toString() + ')'
    }

    return label
  }

  get dataActionsLabel(): string {
    let label = '动作'

    if (this.dataRunSave.actions.length > 0) {
      label += '(' + this.dataRunSave.actions.length.toString() + ')'
    }

    return label
  }

  get dataAssertsLabel(): string {
    let label = '断言'

    if (this.dataRunSave.asserts.length > 0) {
      label += '(' + this.dataRunSave.asserts.length.toString() + ')'
    }

    return label
  }

  get historyActionsLabel(): string {
    let label = '动作'

    if (this.historyRunSave.actions.length > 0) {
      label += '(' + this.historyRunSave.actions.length.toString() + ')'
    }

    return label
  }

  get historyAssertsLabel(): string {
    let label = '断言'

    if (this.historyRunSave.asserts.length > 0) {
      label += '(' + this.historyRunSave.asserts.length.toString() + ')'
    }

    return label
  }

  get respsLabel(): string {
    let label = 'Resps'

    if (this.apiRunSave.respVars.length > 0) {
      label += '(' + this.apiRunSave.respVars.length.toString() + ')'
    }

    return label
  }

  get historyRespsLabel(): string {
    let label = 'Resps'

    if (this.historyRunSave.respVars.length > 0) {
      label += '(' + this.historyRunSave.respVars.length.toString() + ')'
    }

    return label
  }

  get datasLabel(): string {
    let label = '数据列表'
    if (this.sceneSave.dataList.length > 0) {
      label += '(' + this.sceneSave.dataList.length.toString() + ')'
    }
    return label
  }

  get preApisLabel(): string {
    let label = '前置数据'

    if (this.apiRunSave.preApis.length > 0) {
      label += '(' + this.apiRunSave.preApis.length.toString() + ')'
    }

    return label
  }

  get dataPreApisLabel(): string {
    let label = '前置数据'

    if (this.dataRunSave.preApis.length > 0) {
      label += '(' + this.dataRunSave.preApis.length.toString() + ')'
    }

    return label
  }

  get historyPreApisLabel(): string {
    let label = '前置数据'

    if (this.historyRunSave.preApis.length > 0) {
      label += '(' + this.historyRunSave.preApis.length.toString() + ')'
    }

    return label
  }

  get postApisLabel(): string {
    let label = '后置数据'

    if (this.apiRunSave.postApis.length > 0) {
      label += '(' + this.apiRunSave.postApis.length.toString() + ')'
    }

    return label
  }

  get dataPostApisLabel(): string {
    let label = '后置数据'

    if (this.dataRunSave.postApis.length > 0) {
      label += '(' + this.dataRunSave.postApis.length.toString() + ')'
    }

    return label
  }

  get historyPostApisLabel(): string {
    let label = '后置数据'

    if (this.historyRunSave.postApis.length > 0) {
      label += '(' + this.historyRunSave.postApis.length.toString() + ')'
    }

    return label
  }

  get otherConfigLabel(): string {
    let label = '其他'

    if (this.apiRunSave.otherConfigs.length > 0) {
      label += '(' + this.apiRunSave.otherConfigs.length.toString() + ')'
    }

    return label
  }

  get dataOtherConfigLabel(): string {
    let label = '其他'

    if (this.dataRunSave.otherConfigs.length > 0) {
      label += '(' + this.dataRunSave.otherConfigs.length.toString() + ')'
    }

    return label
  }

  get historyOtherConfigLabel(): string {
    let label = '其他'

    if (this.historyRunSave.otherConfigs.length > 0) {
      label += '(' + this.historyRunSave.otherConfigs.length.toString() + ')'
    }

    return label
  }

  get isBody(): boolean {
    let isBody = this.req.method != 'GET' && this.req.method != 'HEAD'
    if (!isBody && this.
        reqMode === 'Body') {
      this.reqMode = 'Header'
    }

    return isBody
  }

  get isDefBodyRaw(): boolean {
    let isDefBodyRaw = false

    if (this.apiDefSave.bodyMode === "application/json") {
      isDefBodyRaw = false
    } else if (this.apiDefSave.bodyMode === "multipart/form-data") {
      isDefBodyRaw = false
    } else if (this.apiDefSave.bodyMode === "application/x-www-form-urlencoded") {
      isDefBodyRaw = false
    } else {
      isDefBodyRaw = true
    }

    return isDefBodyRaw
  }

  get isRunBodyRaw(): boolean {
    // let isRunBodyRaw = this.apiRunSave.bodyMode == 'raw'
    // if (isRunBodyRaw) {
    //   this.apiRunSave.bodyMode = this.rawRunContentType
    // }
    let isRunBodyRaw = false

    if (this.apiRunSave.bodyMode === "application/json") {
      isRunBodyRaw = false
    } else if (this.apiRunSave.bodyMode === "multipart/form-data") {
      isRunBodyRaw = false
    } else if (this.apiRunSave.bodyMode === "application/x-www-form-urlencoded") {
      isRunBodyRaw = false
    } else {
      isRunBodyRaw = true
    }

    return isRunBodyRaw
  }

  get isDataBodyRaw(): boolean {
    // let isDataBodyRaw = this.dataRunSave.bodyMode == 'raw'
    //
    // if (isDataBodyRaw) {
    //   this.dataRunSave.bodyMode = this.rawDefContentType
    // }

    let isDataBodyRaw = false

    if (this.dataRunSave.bodyMode === "application/json") {
      isDataBodyRaw = false
    } else if (this.dataRunSave.bodyMode === "multipart/form-data") {
      isDataBodyRaw = false
    } else if (this.dataRunSave.bodyMode === "application/x-www-form-urlencoded") {
      isDataBodyRaw = false
    } else {
      isDataBodyRaw = true
    }

    return isDataBodyRaw
  }

  get isHistoryBodyRaw(): boolean {
    // let isHistoryBodyRaw = this.historyRunSave.bodyMode == 'raw'
    // if (isHistoryBodyRaw) {
    //   this.historyRunSave.bodyMode = this.rawHistoryContentType
    // }

    let isHistoryBodyRaw = false

    if (this.historyRunSave.bodyMode === "application/json") {
      isHistoryBodyRaw = false
    } else if (this.historyRunSave.bodyMode === "multipart/form-data") {
      isHistoryBodyRaw = false
    } else if (this.historyRunSave.bodyMode === "application/x-www-form-urlencoded") {
      isHistoryBodyRaw = false
    } else {
      isHistoryBodyRaw = true
    }

    return isHistoryBodyRaw
  }

  created() {
    this.onLocale()
  }

  @Watch('$i18n.locale')
  onLocale() {
    this.onTabLabelChange()
  }

  mounted() {
    setTimeout(() => {
      this.isAdvanced = false
      this.onParamsShow()

    }, 1)
    this.getAppList();
    this.getDataList();
    this.getMenu();
    this.getDataMenu();
    this.getSceneMenu();
    this.getHistoryMenu();
    this.getSceneHistoryMenu();
    this.getEnvList();
    this.getDataFileList();
    this.getPlaybookList();
  }

  setDomainMode(value) {
    this.domainMode = value
  }

  onPreview() {
    this.isPreview = !this.isPreview

    if (this.isPreview) {
      this.previewIcon = 'ios-arrow-up'
      this.previewType = 'success'
    } else {
      this.previewIcon = 'ios-arrow-down'
      this.previewType = 'primary'
    }
  }

  onParamsShow() {
    this.showParams = !this.showParams
  }

  onTabLabelChange() {
    let label = this.req.url
    if (label.trim() == '') {
      label = this.$t('tabTitle.newTab').toString()
    }

    this.tab.label = label.length > 18 ? label.substr(0, 18) : label
  }

  onAdvanced() {
    this.isAdvanced = !this.isAdvanced

    if (this.isAdvanced) {
      this.advancedIcon = 'ios-arrow-up'
      this.isPreview = false
    } else {
      this.advancedIcon = 'ios-arrow-down'
      this.isPreview = true
    }
  }

  async onDefSave() {
    this.isSending = true
    let data = new URLSearchParams()
    data.append('app', JSON.stringify(this.apiDefSave.app))
    data.append('module', JSON.stringify(this.apiDefSave.module))
    data.append('apiDesc', JSON.stringify(this.apiDefSave.apiDesc))
    data.append('method', JSON.stringify(this.apiDefSave.method))
    data.append('path', JSON.stringify(this.apiDefSave.path))
    data.append('prefix', JSON.stringify(this.apiDefSave.prefix))
    data.append('pathVars', JSON.stringify(this.apiDefSave.pathVars))
    data.append('queryVars', JSON.stringify(this.apiDefSave.queryVars))
    data.append('bodyVars', JSON.stringify(this.apiDefSave.bodyVars))
    data.append('headerVars', JSON.stringify(this.apiDefSave.headerVars))
    data.append('respVars', JSON.stringify(this.apiDefSave.respVars))
    data.append('bodyMode', JSON.stringify(this.apiDefSave.bodyMode))
    data.append('bodyStr', JSON.stringify(this.apiDefSave.bodyStr))

    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/apiDefSave', data)
    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 5,
        content: result.msg + '(' + result.code.toString() + ')'
      })
    }

    this.isSending = false
    this.getMenu()
  }

  async onApiRunSave() {
    this.isSending = true
    let data = new URLSearchParams()
    data.append('app', JSON.stringify(this.apiRunSave.app))
    data.append('module', JSON.stringify(this.apiRunSave.module))
    data.append('apiDesc', JSON.stringify(this.apiRunSave.apiDesc))
    data.append('dataDesc', JSON.stringify(this.apiRunSave.dataDesc))
    data.append('method', JSON.stringify(this.apiRunSave.method))
    data.append('path', JSON.stringify(this.apiRunSave.path))
    data.append('prefix', JSON.stringify(this.apiRunSave.prefix))
    data.append('pathVars', JSON.stringify(this.apiRunSave.pathVars))
    data.append('queryVars', JSON.stringify(this.apiRunSave.queryVars))
    data.append('bodyVars', JSON.stringify(this.apiRunSave.bodyVars))
    data.append('bodyStr', JSON.stringify(this.apiRunSave.bodyStr))
    data.append('headerVars', JSON.stringify(this.apiRunSave.headerVars))
    data.append('respVars', JSON.stringify(this.apiRunSave.respVars))
    data.append('actions', JSON.stringify(this.apiRunSave.actions))
    data.append('asserts', JSON.stringify(this.apiRunSave.asserts))
    data.append('preApis', JSON.stringify(this.apiRunSave.preApis))
    data.append('postApis', JSON.stringify(this.apiRunSave.postApis))
    data.append('product', JSON.stringify(this.apiRunSave.product))
    data.append('bodyMode', JSON.stringify(this.apiRunSave.bodyMode))

    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/apiDataSave', data)
    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 5,
        content: result.msg + '(' + result.code.toString() + ')'
      })
    }
    this.isSending = false
    this.getDataMenu();
  }

  async onDataSave() {
    this.isSending = true
    let data = new URLSearchParams()
    data.append('app', JSON.stringify(this.dataRunSave.app))
    data.append('module', JSON.stringify(this.dataRunSave.module))
    data.append('apiDesc', JSON.stringify(this.dataRunSave.apiDesc))
    data.append('dataDesc', JSON.stringify(this.dataRunSave.dataDesc))
    data.append('method', JSON.stringify(this.dataRunSave.method))
    data.append('path', JSON.stringify(this.dataRunSave.path))
    data.append('prefix', JSON.stringify(this.dataRunSave.prefix))
    data.append('pathVars', JSON.stringify(this.dataRunSave.pathVars))
    data.append('queryVars', JSON.stringify(this.dataRunSave.queryVars))
    data.append('bodyVars', JSON.stringify(this.dataRunSave.bodyVars))
    data.append('bodyStr', JSON.stringify(this.dataRunSave.bodyStr))
    data.append('headerVars', JSON.stringify(this.dataRunSave.headerVars))
    data.append('respVars', JSON.stringify(this.dataRunSave.respVars))
    data.append('actions', JSON.stringify(this.dataRunSave.actions))
    data.append('asserts', JSON.stringify(this.dataRunSave.asserts))
    data.append('preApis', JSON.stringify(this.dataRunSave.preApis))
    data.append('postApis', JSON.stringify(this.dataRunSave.postApis))
    data.append('product', JSON.stringify(this.dataRunSave.product))
    data.append('otherConfig', JSON.stringify(this.dataRunSave.otherConfigs))
    data.append('bodyMode', JSON.stringify(this.dataRunSave.bodyMode))

    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/apiDataSave', data)
    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 5,
        content: result.msg + '(' + result.code.toString() + ')'
      })
    }
    this.isSending = false
    this.getDataMenu();
  }

  async onHistorySave() {
    this.isSending = true
    let data = new URLSearchParams()
    data.append('app', JSON.stringify(this.historyRunSave.app))
    data.append('module', JSON.stringify(this.historyRunSave.module))
    data.append('apiDesc', JSON.stringify(this.historyRunSave.apiDesc))
    data.append('dataDesc', JSON.stringify(this.historyRunSave.dataDesc))
    data.append('method', JSON.stringify(this.historyRunSave.method))
    data.append('path', JSON.stringify(this.historyRunSave.path))
    data.append('host', JSON.stringify(this.historyRunSave.host))
    data.append('prototype', JSON.stringify(this.historyRunSave.prototype))
    data.append('prefix', JSON.stringify(this.historyRunSave.prefix))
    data.append('pathVars', JSON.stringify(this.historyRunSave.pathVars))
    data.append('queryVars', JSON.stringify(this.historyRunSave.queryVars))
    data.append('bodyVars', JSON.stringify(this.historyRunSave.bodyVars))
    data.append('bodyStr', JSON.stringify(this.historyRunSave.bodyStr))
    data.append('headerVars', JSON.stringify(this.historyRunSave.headerVars))
    data.append('respVars', JSON.stringify(this.historyRunSave.respVars))
    data.append('actions', JSON.stringify(this.historyRunSave.actions))
    data.append('asserts', JSON.stringify(this.historyRunSave.asserts))
    data.append('preApis', JSON.stringify(this.historyRunSave.preApis))
    data.append('postApis', JSON.stringify(this.historyRunSave.postApis))
    data.append('product', JSON.stringify(this.historyRunSave.product))
    data.append('otherConfig', JSON.stringify(this.historyRunSave.otherConfigs))
    data.append('bodyMode', JSON.stringify(this.historyRunSave.bodyMode))

    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/historySave', data)
    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 5,
        content: result.msg + '(' + result.code.toString() + ')'
      })
    }

    this.isSending = false
    this.getHistoryMenu()
  }

  async onSceneSave() {
    this.isSending = true
    let data = new URLSearchParams()
    data.append('product', JSON.stringify(this.sceneSave.product))
    data.append('name', JSON.stringify(this.sceneSave.name))
    data.append('dataList', JSON.stringify(this.sceneSave.dataList))
    data.append('type', JSON.stringify(this.sceneSave.type))
    data.append('runNum', JSON.stringify(this.sceneSave.runNum))

    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/sceneSave', data)
    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 5,
        content: result.msg + '(' + result.code.toString() + ')'
      })
    }

    this.isSending = false
    this.getSceneMenu()
  }

  async onSceneHistorySave() {
    this.isSending = true
    let data = new URLSearchParams()
    data.append('product', JSON.stringify(this.sceneHistorySave.product))
    data.append('name', JSON.stringify(this.sceneHistorySave.name))
    data.append('dataList', JSON.stringify(this.sceneHistorySave.dataList))
    data.append('type', JSON.stringify(this.sceneHistorySave.type))
    data.append('runNum', JSON.stringify(this.sceneHistorySave.runNum))

    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/sceneSave', data)
    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 5,
        content: result.msg + '(' + result.code.toString() + ')'
      })
    }

    this.isSending = false
    this.getSceneMenu()
  }

  onMock() {

  }

  onSave() {
    this.isSaved = !this.isSaved

    if (this.isSaved) {
      this.saveIcon = 'save'
      this.isPreview = false
    } else {
      this.saveIcon = 'save'
      this.isPreview = true
    }
  }

  async getAppList() {
    if (this.appOptions.length===0) {
      let result = await API.get<Req.ResponseModel[]>('/appList')
      if (result.data && result.data.length > 0) {
        this.appOptions = []
        _.forEach(result.data, v => {
          this.appOptions = this.appOptions.concat(v+'')
        })
      }
    }
  }

  async getPlaybookList() {
    if (this.playbookOptions.length===0) {
      let result = await API.get<Req.ResponseModel[]>('/sceneList')
      if (result.data && result.data.length > 0) {
        this.playbookOptions = []
        _.forEach(result.data, v => {
          this.playbookOptions = this.playbookOptions.concat(v+'')
        })
      }
    }
  }

  async getMenu() {
    if (this.domainMode === "run") {
      let result = await API.get<Req.ResponseModel[]>('/allMenu?app='+this.apiRunSave.app)
      if (result.data && result.data.length > 0) {
        this.menuData = []
        _.forEach(result.data, v => {
          this.menuData = this.menuData.concat(v)
        })
      }
    } else {
      let result = await API.get<Req.ResponseModel[]>('/allMenu?app='+this.apiDefSave.app)
      if (result.data && result.data.length > 0) {
        this.menuData = []
        _.forEach(result.data, v => {
          this.menuData = this.menuData.concat(v)
        })
      }
    }
  }

  async getDataMenu() {
    let result = await API.get<Req.ResponseModel[]>('/dataMenu?app='+this.dataRunSave.app)
    if (result.data && result.data.length > 0) {
      this.dataMenu = []
      _.forEach(result.data, v => {
        this.dataMenu = this.dataMenu.concat(v)
      })
    }
  }

  async getSceneMenu() {
    let result = await API.get<Req.ResponseModel[]>('/sceneMenu?product='+this.sceneSave.product)
    if (result.data && result.data.length > 0) {
      this.sceneMenu = []
      _.forEach(result.data, v => {
        this.sceneMenu = this.sceneMenu.concat(v)
      })
    }
  }

  async getHistoryMenu() {
    let result = await API.get<Req.ResponseModel[]>('/historyMenu?dateName='+this.selectDate)
    if (result.data && result.data.length > 0) {
      this.historyMenu = []
      _.forEach(result.data, v => {
        this.historyMenu = this.historyMenu.concat(v)
      })
    }
  }

  async getSceneHistoryMenu() {
    let result = await API.get<Req.ResponseModel[]>('/sceneHistoryMenu?dateName='+this.sceneSelectDate)
    if (result.data && result.data.length > 0) {
      this.sceneHistoryMenu = []
      _.forEach(result.data, v => {
        this.sceneHistoryMenu = this.sceneHistoryMenu.concat(v)
      })
    }
  }

  async getDataList() {
    if (this.dataDescOptions.length===0) {
      let result = await API.get<Req.ResponseModel[]>('/dataList')
      if (result.data && result.data.length > 0) {
        this.dataDescOptions = []
        _.forEach(result.data, v => {
          this.dataDescOptions = this.dataDescOptions.concat(v+'')
        })
      }
    }
  }

  async getAppData() {
    let appName = ""
    if (this.domainMode==="def") {
      if (this.apiDefSave.app) {
        appName = this.apiDefSave.app
      }
      this.apiDefSave.method = ""
      this.apiDefSave.path = ""
      this.apiDefSave.module = ""
      this.apiDefSave.apiDesc = ""
      this.apiDefSave.prefix = ""
    } else if (this.domainMode==="run") {
      if (this.apiRunSave.app) {
        appName = this.apiRunSave.app
      }
      this.apiRunSave.method = ""
      this.apiRunSave.path = ""
      this.apiRunSave.module = ""
      this.apiRunSave.apiDesc = ""
      this.apiRunSave.prefix = ""
    }

    if (appName.length>0) {
      let result = await API.get<Req.ResponseModel[]>('/appList/'+appName)
      if (result.data) {
        if (this.domainMode==="def") {
          if (result.data["methods"]!=null) {
            this.defMethodOptions = result.data["methods"]
          } else {
            this.defMethodOptions = []
          }
          if (result.data["modules"]!=null) {
            this.defModuleOptions = result.data["modules"]
          } else {
            this.defModuleOptions = []
          }
          if (result.data["apis"]!=null) {
            this.defApiOptions = result.data["apis"]
          } else {
            this.defApiOptions = []
          }
          if (result.data["apisDesc"]!=null) {
            this.defApiDescOptions = result.data["apisDesc"]
          } else {
            this.defApiDescOptions = []
          }
          this.apiDefSave.prefix = result.data["prefix"]
        } else if (this.domainMode==="run") {
          if (result.data["methods"]!=null) {
            this.runMethodOptions = result.data["methods"]
          } else {
            this.runMethodOptions = []
          }
          if (result.data["modules"]!=null) {
            this.runModuleOptions = result.data["modules"]
          } else {
            this.runModuleOptions = []
          }
          if (result.data["apisDesc"]!=null) {
            this.runApiDescOptions = result.data["apisDesc"]
          } else {
            this.runApiDescOptions = []
          }
          if (result.data["datasDesc"]!=null) {
            this.dataDescOptions = result.data["datasDesc"]
          } else {
            this.dataDescOptions = []
          }
          this.apiRunSave.prefix = result.data["prefix"]
        }
      }
    }
  }

  async getDataApp() {
    let appName = ""
    appName = this.dataRunSave.app
    if (appName.length>0) {
      let result = await API.get<Req.ResponseModel[]>('/appList/'+appName)
      if (result.data) {
        this.dataRunSave.prefix = result.data["prefix"]
      }
    }
  }

  async getModuleData() {
    let appName = ""
    let module = ""
    if (this.domainMode==="def") {
      if (this.apiDefSave.app) {
        appName = this.apiDefSave.app
      }
      if (this.apiDefSave.module) {
        module = this.apiDefSave.module
      }
    } else if (this.domainMode==="run") {
      if (this.apiRunSave.app) {
        appName = this.apiRunSave.app
      }
      if (this.apiRunSave.module) {
        module = this.apiRunSave.module
      }
    }

    if (appName.length>0 && module.length>0) {
      let result = await API.get<Req.ResponseModel[]>('/appList/'+appName+"?module="+module+"&mode="+this.domainMode)
      if (result.data) {
        if (this.domainMode==="def") {
          this.defMethodOptions = result.data["methods"]
          this.defApiOptions = result.data["apis"]
          this.defApiDescOptions = result.data["apisDesc"]
        } else if (this.domainMode==="run") {
          this.runMethodOptions = result.data["methods"]
          this.runApiOptions = result.data["apis"]
          this.runApiDescOptions = result.data["apisDesc"]
        }
      }
    }
  }

  async getMethodData() {
    let appName = ""
    let method = ""

    if (this.domainMode==="def") {
      if (this.apiDefSave.app) {
        appName = this.apiDefSave.app
      }
      if (this.apiDefSave.method) {
        method = this.apiDefSave.method
      }
    } else if (this.domainMode==="run") {
      if (this.apiRunSave.app) {
        appName = this.apiRunSave.app
      }
      if (this.apiRunSave.module) {
        method = this.apiRunSave.method
      }
    }

    if (appName.length>0 && method.length>0) {
      let result = await API.get<Req.ResponseModel[]>('/appList/'+appName+"?method="+method+"&mode="+this.domainMode)
      if (result.data) {
        if (this.domainMode==="def") {
          this.defApiOptions = result.data["apis"]
          this.defApiDescOptions = result.data["apisDesc"]
          this.defModuleOptions = result.data["modules"]
        } else  if (this.domainMode==="run") {
          this.runApiOptions = result.data["apis"]
          this.runApiDescOptions = result.data["apisDesc"]
          this.runModuleOptions = result.data["modules"]
        }
      }
    }
  }

  selectApiDetail(val) {
    this.apiDefSave.apiDesc = val
    this.getApiDescData()
  }

  selectDataDetail(val) {
    this.dataRunSave.app = val[0]
    this.getApiDataDetail()
  }

  selectRunApiDesc(val) {
    this.apiRunSave.apiDesc = val
    this.getApiDescData()
  }

  selectDefModule(val) {
    if (val.length === 1 ) {
      this.apiDefSave.app = val[0]
      this.getMenu()
    } else if (val.length === 2 ) {
      this.apiDefSave.app = val[0]
      this.apiDefSave.module = val[1]
    }
  }

  selectData(val) {
    this.dataRunSave.dataDesc = val
    this.getDataDetail()
  }

  getSceneByName(val) {
    this.sceneSave.name = val
    this.getPlaybookDetail()
  }

  selectScene(val) {
    this.sceneSave.name = val
    this.getSceneDetail()
  }

  selectHistory(val) {
    this.historyRunSave.fileName = val
    this.getHistoryDetail()
  }

  selectSceneHistory(val) {
    this.sceneHistorySave.name = val
    this.getSceneHistoryDetail()
  }

  selectSceneDetail(val) {
    if (val.length > 0 ) {
      this.sceneSave.product = val[0]
    }
  }

  selectRunModule(val) {
    if (val.length === 1 ) {
      this.apiRunSave.app = val[0]
      this.getMenu()
    } else if (val.length === 2 ) {
      this.apiRunSave.app = val[0]
      this.apiRunSave.module = val[1]
    }
  }

  selectAppData(val) {
    this.dataRunSave.app = val[0]
    this.getDataMenu()
  }

  selectProductScene(val) {
    this.sceneSave.product = val[0]
    this.getSceneMenu()
  }

  selectHistoryDate(val) {
    this.selectDate = val[0]
    this.getHistoryMenu()
  }

  selectSceneHistoryDate(val) {
    this.sceneSelectDate = val[0]
    this.getSceneHistoryMenu()
  }

  async getApiDescData() {
    let appName = ""
    let module = ""
    let apiDesc = ""
    if (this.domainMode==="def") {
      if (this.apiDefSave.app) {
        appName = this.apiDefSave.app
      }
      if (this.apiDefSave.module) {
        module = this.apiDefSave.module
      }
      if (this.apiDefSave.apiDesc) {
        apiDesc = this.apiDefSave.apiDesc
      }
    } else if (this.domainMode==="run") {
      if (this.apiRunSave.app) {
        appName = this.apiRunSave.app
      }
      if (this.apiRunSave.module) {
        module = this.apiRunSave.module
      }
      if (this.apiRunSave.apiDesc) {
        apiDesc = this.apiRunSave.apiDesc
      }
    }
    if (appName.length>0 && module.length > 0 && apiDesc.length > 0){
      let result = await API.get<Req.ResponseModel[]>('/appList/'+appName+"?module="+module+"&apiDesc="+apiDesc+"&mode="+this.domainMode)
      if (result.data) {
        if (this.domainMode==="def") {
          this.apiDefSave.method = result.data["method"]
          this.apiDefSave.path = result.data["path"]
          this.apiDefSave.module = result.data["module"]
          this.apiDefSave.prefix = result.data["prefix"]
          this.apiDefSave.apiDesc = result.data["apiDesc"]
          this.apiDefSave.bodyMode = result.data["bodyMode"]

          this.apiDefSave.pathVars = []
          if (result.data["pathVars"]) {
            _.forEach(result.data["pathVars"], v => {
              this.apiDefSave.pathVars = this.apiDefSave.pathVars.concat(v)
            })
          }
          this.apiDefSave.queryVars = []
          if (result.data["queryVars"]) {
            _.forEach(result.data["queryVars"], v => {
              this.apiDefSave.queryVars = this.apiDefSave.queryVars.concat(v)
            })
          }
          this.apiDefSave.bodyVars = []
          if (result.data["bodyVars"]) {
            _.forEach(result.data["bodyVars"], v => {
              this.apiDefSave.bodyVars = this.apiDefSave.bodyVars.concat(v)
            })
          }
          this.apiDefSave.headerVars = []
          if (result.data["headerVars"]) {
            _.forEach(result.data["headerVars"], v => {
              this.apiDefSave.headerVars = this.apiDefSave.headerVars.concat(v)
            })
          }
          this.apiDefSave.respVars = []
          if (result.data["respVars"]) {
            _.forEach(result.data["respVars"], v => {
              this.apiDefSave.respVars = this.apiDefSave.respVars.concat(v)
            })
          }
        } else if (this.domainMode==="run") {
          this.apiRunSave.method = result.data["method"]
          this.apiRunSave.path = result.data["path"]
          this.apiRunSave.module = result.data["module"]
          this.apiRunSave.prefix = result.data["prefix"]
          this.apiRunSave.dataDesc = ''

          if (result.data["datasDesc"]!=null) {
            this.dataDescOptions = result.data["datasDesc"]
          } else {
            this.dataDescOptions = []
          }

          this.apiRunSave.pathVars = []
          if (result.data["pathVars"]) {
            _.forEach(result.data["pathVars"], v => {
              this.apiRunSave.pathVars = this.apiRunSave.pathVars.concat(v)
            })
          }
          this.apiRunSave.queryVars = []
          if (result.data["queryVars"]) {
            _.forEach(result.data["queryVars"], v => {
              this.apiRunSave.queryVars = this.apiRunSave.queryVars.concat(v)
            })
          }
          this.apiRunSave.bodyVars = []
          if (result.data["bodyVars"]) {
            _.forEach(result.data["bodyVars"], v => {
              this.apiRunSave.bodyVars = this.apiRunSave.bodyVars.concat(v)
            })
          }
          this.apiRunSave.headerVars = []
          if (result.data["headerVars"]) {
            _.forEach(result.data["headerVars"], v => {
              this.apiRunSave.headerVars = this.apiRunSave.headerVars.concat(v)
            })
          }
          this.apiRunSave.respVars = []
          if (result.data["respVars"]) {
            _.forEach(result.data["respVars"], v => {
              this.apiRunSave.respVars = this.apiRunSave.respVars.concat(v)
            })
          }
          this.apiRunSave.actions = []
          if (result.data["actions"]) {
            _.forEach(result.data["actions"], v => {
              this.apiRunSave.actions = this.apiRunSave.actions.concat(v)
            })
          }
          this.apiRunSave.asserts = []
          if (result.data["asserts"]) {
            _.forEach(result.data["asserts"], v => {
              this.apiRunSave.asserts = this.apiRunSave.asserts.concat(v)
            })
          }

          this.apiRunSave.otherConfigs = []
          if (result.data["otherConfig"]) {
            _.forEach(result.data["otherConfig"], v => {
              this.apiRunSave.otherConfigs = this.apiRunSave.otherConfigs.concat(v)
            })
          }
        }
      }
    }
  }

  async getSceneDetail() {
    let product = this.sceneSave.product
    let name = this.sceneSave.name

    let result = await API.get< Req.RelatedApiListModel[]>('/sceneList/'+name+"?product="+product)
    if (result.data) {
      this.sceneSave.dataList = []
      _.forEach(result.data['dataList'], v => {
        this.sceneSave.dataList = this.sceneSave.dataList.concat(v)
      })
      this.sceneSave.type = result.data['type']
      this.sceneSave.runNum = result.data['runNum']
    }
  }

  async getPlaybookDetail() {
    let name = this.sceneSave.name
    let result = await API.get< Req.RelatedApiListModel[]>('/sceneList/'+name)
    if (result.code === 200) {
      if (result.data) {
        this.sceneSave.dataList = []
        _.forEach(result.data['dataList'], v => {
          this.sceneSave.dataList = this.sceneSave.dataList.concat(v)
        })
        this.sceneSave.type = result.data['type']
        this.sceneSave.runNum = result.data['runNum']
        this.sceneSave.product = result.data['product']
      }
    }
  }


  async getHistoryDetail() {
    let result = await API.get<Req.ResponseModel[]>('/historyList?fileName='+this.historyRunSave.fileName)
    if (result.code == 200) {
      if (result.data) {
        this.historyRunSave.method = result.data["method"]
        this.historyRunSave.path = result.data["path"]
        this.historyRunSave.module = result.data["module"]
        this.historyRunSave.prefix = result.data["prefix"]
        this.historyRunSave.dataDesc = result.data["dataDesc"]
        this.historyRunSave.app = result.data["app"]
        this.historyRunSave.apiDesc = result.data["apiDesc"]
        this.historyRunSave.prototype = result.data["prototype"]
        this.historyRunSave.host = result.data["host"]
        this.historyRunSave.output = result.data["output"]
        this.historyRunSave.bodyMode = result.data["bodyMode"]

        this.historyRunSave.pathVars = []
        if (result.data["pathVars"]) {
          _.forEach(result.data["pathVars"], v => {
            this.historyRunSave.pathVars = this.historyRunSave.pathVars.concat(v)
          })
        }

        this.historyRunSave.queryVars = []
        if (result.data["queryVars"]) {
          _.forEach(result.data["queryVars"], v => {
            this.historyRunSave.queryVars = this.historyRunSave.queryVars.concat(v)
          })
        }

        this.historyRunSave.bodyVars = []
        if (result.data["bodyVars"]) {
          _.forEach(result.data["bodyVars"], v => {
            this.historyRunSave.bodyVars = this.historyRunSave.bodyVars.concat(v)
          })
        }

        this.historyRunSave.headerVars = []
        if (result.data["headerVars"]) {
          _.forEach(result.data["headerVars"], v => {
            this.historyRunSave.headerVars = this.historyRunSave.headerVars.concat(v)
          })
        }

        this.historyRunSave.respVars = []
        if (result.data["respVars"]) {
          _.forEach(result.data["respVars"], v => {
            this.historyRunSave.respVars = this.historyRunSave.respVars.concat(v)
          })
        }

        this.historyRunSave.actions = []
        if (result.data["actions"]) {
          _.forEach(result.data["actions"], v => {
            this.historyRunSave.actions = this.historyRunSave.actions.concat(v)
          })
        }

        this.historyRunSave.asserts = []
        if (result.data["asserts"]) {
          _.forEach(result.data["asserts"], v => {
            this.historyRunSave.asserts = this.historyRunSave.asserts.concat(v)
          })
        }

        this.apiRunSave.otherConfigs = []
        if (result.data["otherConfig"]) {
          _.forEach(result.data["otherConfig"], v => {
            this.historyRunSave.otherConfigs = this.historyRunSave.otherConfigs.concat(v)
          })
        }

        this.historyModeReqDataRespList.url = result.data["url"]
        this.historyModeReqDataRespList.response = result.data["response"]
        this.historyModeReqDataRespList.failReason = result.data["failReason"]
        this.historyModeReqDataRespList.testResult = result.data["testResult"]
        this.historyModeReqDataRespList.request = result.data["request"]
        this.historyModeReqDataRespList.header = result.data["header"]
        this.historyModeReqDataRespList.output = result.data["output"]

      }
    } else {
      this.$Message.error({
        duration: 10,
        content: result.msg + '(' + result.code.toString() + ')'
      })
    }
  }

  async getSceneHistoryDetail() {
    let result = await API.get<Req.ResponseModel[]>('/sceneHistoryList?name='+this.sceneHistorySave.name)
    if (result.code == 200) {
      if (result.data) {
        this.sceneHistorySave.name = result.data["name"]
        this.sceneHistorySave.product = result.data["product"]
        this.sceneHistorySave.type = result.data["type"]
        this.sceneHistorySave.runNum = result.data["runNum"]
        this.sceneHistorySave.dataList = result.data["dataList"]
        this.sceneHistorySave.lastFile = result.data["lastDataFile"]
        this.sceneHistoryRespList.lastDataFile = result.data["lastDataFile"]
        this.sceneHistoryRespList.failReason = result.data["failReason"]
        this.sceneHistoryRespList.testResult = result.data["testResult"]
      }
    } else {
      this.$Message.error({
        duration: 10,
        content: result.msg + '(' + result.code.toString() + ')'
      })
    }
  }

  async getApiDescDataList() {
    let appName = ""
    let module = ""
    let apiDesc = ""
    if (this.domainMode==="def") {
      if (this.apiDefSave.app) {
        appName = this.apiDefSave.app
      }
      if (this.apiDefSave.module) {
        module = this.apiDefSave.module
      }
      if (this.apiDefSave.apiDesc) {
        apiDesc = this.apiDefSave.apiDesc
      }
    } else if (this.domainMode==="run") {
      if (this.apiRunSave.app) {
        appName = this.apiRunSave.app
      }
      if (this.apiRunSave.module) {
        module = this.apiRunSave.module
      }
      if (this.apiRunSave.apiDesc) {
        apiDesc = this.apiRunSave.apiDesc
      }
    }

    if (appName.length>0 && module.length > 0 && apiDesc.length > 0){
      let result = await API.get<Req.ResponseModel[]>('/appList/'+appName+"?module="+module+"&apiDesc="+appName+"&mode="+this.domainMode)
      if (result.data) {
        if (this.domainMode==="def") {
          this.apiDefSave.method = result.data["method"]
          this.apiDefSave.path = result.data["path"]
          this.apiDefSave.pathVars = []
          if (result.data["pathVars"]) {
            _.forEach(result.data["pathVars"], v => {
              // let varModel = v
              this.apiDefSave.pathVars = this.apiDefSave.pathVars.concat(v)
            })
          }
          this.apiDefSave.queryVars = []
          if (result.data["queryVars"]) {
            _.forEach(result.data["queryVars"], v => {
              // let varModel = v
              this.apiDefSave.queryVars = this.apiDefSave.queryVars.concat(v)
            })
          }
          this.apiDefSave.bodyVars = []
          if (result.data["bodyVars"]) {
            _.forEach(result.data["bodyVars"], v => {
              // let varModel = v
              this.apiDefSave.bodyVars = this.apiDefSave.bodyVars.concat(v)
            })
          }
          this.apiDefSave.headerVars = []
          if (result.data["headerVars"]) {
            _.forEach(result.data["headerVars"], v => {
              // let varModel = v
              this.apiDefSave.headerVars = this.apiDefSave.headerVars.concat(v)
            })
          }
          this.apiDefSave.respVars = []
          if (result.data["respVars"]) {
            _.forEach(result.data["respVars"], v => {
              // let varModel = v
              this.apiDefSave.respVars = this.apiDefSave.respVars.concat(v)
            })
          }
        } else if (this.domainMode==="run") {
          this.apiRunSave.method = result.data["method"]
          this.apiRunSave.path = result.data["path"]
          this.apiRunSave.pathVars = []
          if (result.data["pathVars"]) {
            _.forEach(result.data["pathVars"], v => {
              // let varModel = v
              this.apiRunSave.pathVars = this.apiRunSave.pathVars.concat(v)
            })
          }
          this.apiRunSave.queryVars = []
          if (result.data["queryVars"]) {
            _.forEach(result.data["queryVars"], v => {
              // let varModel = v
              this.apiRunSave.queryVars = this.apiRunSave.queryVars.concat(v)
            })
          }
          this.apiRunSave.bodyVars = []
          if (result.data["bodyVars"]) {
            _.forEach(result.data["bodyVars"], v => {
              // let varModel = v
              this.apiRunSave.bodyVars = this.apiRunSave.bodyVars.concat(v)
            })
          }
          this.apiRunSave.headerVars = []
          if (result.data["headerVars"]) {
            _.forEach(result.data["headerVars"], v => {
              // let varModel = v
              this.apiRunSave.headerVars = this.apiRunSave.headerVars.concat(v)
            })
          }
          this.apiRunSave.respVars = []
          if (result.data["respVars"]) {
            _.forEach(result.data["respVars"], v => {
              // let varModel = v
              this.apiRunSave.respVars = this.apiRunSave.respVars.concat(v)
            })
          }

          this.apiRunSave.otherConfigs = []
          if (result.data["otherConfig"]) {
            _.forEach(result.data["otherConfig"], v => {
              this.apiRunSave.otherConfigs = this.apiRunSave.otherConfigs.concat(v)
            })
          }
        }
      }
    }
  }

  async getDataDetail() {
    let appName = this.dataRunSave.app
    let module = this.dataRunSave.module
    let apiDesc = this.dataRunSave.apiDesc
    let dataDesc = this.dataRunSave.dataDesc
    let method = this.dataRunSave.method
    let path = this.dataRunSave.path

    this.reqDataRespList.response = ""

    if (!dataDesc) {
      return
    }
    if (dataDesc.length > 0){
      let result
      if (appName.length>0) {
        result = await API.get<Req.ResponseModel[]>('/appList/'+appName+"?module="+module+"&apiDesc="+apiDesc+"&dataDesc="+dataDesc+"&method="+method+"&path="+path+"&mode=run")
      } else {
        result = await API.get<Req.ResponseModel[]>('/dataList'+"?dataDesc="+dataDesc)
      }
      if (result.data) {
        if (result.data["prefix"].length > 0 ) {
          this.dataRunSave.prefix = result.data["prefix"]
        }

        if (this.dataRunSave.app.length===0) {
          this.dataRunSave.app = result.data["app"]
        }
        this.dataRunSave.method = result.data["method"]
        if (this.runMethodOptions.indexOf(this.dataRunSave.method)<0) {
          this.runMethodOptions.push(this.dataRunSave.method)
        }

        this.dataRunSave.path = result.data["path"]
        if (this.runApiOptions.indexOf(this.dataRunSave.path)<0) {
          this.runApiOptions.push(this.dataRunSave.path)
        }

        this.dataRunSave.apiDesc = result.data["apiDesc"]
        if (this.runApiDescOptions.indexOf(this.dataRunSave.apiDesc)<0) {
          this.runApiDescOptions.push(this.dataRunSave.apiDesc)
        }

        this.dataRunSave.module = result.data["module"]
        if (this.runModuleOptions.indexOf(this.dataRunSave.module)<0) {
          this.runModuleOptions.push(this.dataRunSave.module)
        }

        this.dataRunSave.pathVars = []
        if (result.data["pathVars"]) {
          _.forEach(result.data["pathVars"], v => {
            this.dataRunSave.pathVars = this.dataRunSave.pathVars.concat(v)

          })
        }

        this.dataRunSave.queryVars = []
        if (result.data["queryVars"]) {
          _.forEach(result.data["queryVars"], v => {
            this.dataRunSave.queryVars = this.dataRunSave.queryVars.concat(v)
          })
        }

        this.dataRunSave.bodyVars = []
        if (result.data["bodyVars"]) {
          _.forEach(result.data["bodyVars"], v => {
            this.dataRunSave.bodyVars = this.dataRunSave.bodyVars.concat(v)
          })
        }
        this.dataRunSave.bodyMode = result.data["bodyMode"]

        this.dataRunSave.headerVars = []
        if (result.data["headerVars"]) {
          _.forEach(result.data["headerVars"], v => {
            this.dataRunSave.headerVars = this.dataRunSave.headerVars.concat(v)
          })
        }

        this.dataRunSave.respVars = []
        if (result.data["respVars"]) {
          _.forEach(result.data["respVars"], v => {
            this.dataRunSave.respVars = this.dataRunSave.respVars.concat(v)
          })
        }

        this.dataRunSave.actions = []
        if (result.data["actions"]) {
          _.forEach(result.data["actions"], v => {
            this.dataRunSave.actions = this.dataRunSave.actions.concat(v)
          })
        }

        this.dataRunSave.asserts = []
        if (result.data["asserts"]) {
          _.forEach(result.data["asserts"], v => {
            this.dataRunSave.asserts = this.dataRunSave.asserts.concat(v)
          })
        }

        this.dataRunSave.preApis = []
        if (result.data["preApis"]) {
          _.forEach(result.data["preApis"], v => {
            this.dataRunSave.preApis = this.dataRunSave.preApis.concat(v)
          })
        }

        this.dataRunSave.postApis = []
        if (result.data["postApis"]) {
          _.forEach(result.data["postApis"], v => {
            this.dataRunSave.postApis = this.dataRunSave.postApis.concat(v)
          })
        }

        this.dataRunSave.otherConfigs = []
        if (result.data["otherConfig"]) {
          _.forEach(result.data["otherConfig"], v => {
            this.dataRunSave.otherConfigs = this.dataRunSave.otherConfigs.concat(v)
          })
        }
      }
    }
  }

  async getDataByName(val) {
    let dataDesc = val

    this.reqDataRespList.response = ""

    if (!dataDesc) {
      return
    }
    let result
    result = await API.get<Req.ResponseModel[]>('/dataList'+"?dataDesc="+dataDesc)
    if (result.data["path"].length===0) {
      return
    }


    if (result.data["prefix"].length > 0 ) {
      this.dataRunSave.prefix = result.data["prefix"]
    }

    if (this.dataRunSave.app.length===0) {
      this.dataRunSave.app = result.data["app"]
    }

    this.dataRunSave.method = result.data["method"]
    if (this.runMethodOptions.indexOf(this.dataRunSave.method)<0) {
      this.runMethodOptions.push(this.dataRunSave.method)
    }

    this.dataRunSave.path = result.data["app"]
    if (this.appOptions.indexOf(this.dataRunSave.app)<0) {
      this.appOptions.push(this.dataRunSave.path)
    }
    this.dataRunSave.path = result.data["path"]
    if (this.runApiOptions.indexOf(this.dataRunSave.path)<0) {
      this.runApiOptions.push(this.dataRunSave.path)
    }

    this.dataRunSave.apiDesc = result.data["apiDesc"]
    if (this.runApiDescOptions.indexOf(this.dataRunSave.apiDesc)<0) {
      this.runApiDescOptions.push(this.dataRunSave.apiDesc)
    }

    this.dataRunSave.module = result.data["module"]
    if (this.runModuleOptions.indexOf(this.dataRunSave.module)<0) {
      this.runModuleOptions.push(this.dataRunSave.module)
    }

    this.dataRunSave.pathVars = []
    if (result.data["pathVars"]) {
      _.forEach(result.data["pathVars"], v => {
        this.dataRunSave.pathVars = this.dataRunSave.pathVars.concat(v)

      })
    }

    this.dataRunSave.queryVars = []
    if (result.data["queryVars"]) {
      _.forEach(result.data["queryVars"], v => {
        this.dataRunSave.queryVars = this.dataRunSave.queryVars.concat(v)
      })
    }

    this.dataRunSave.bodyVars = []
    if (result.data["bodyVars"]) {
      _.forEach(result.data["bodyVars"], v => {
        this.dataRunSave.bodyVars = this.dataRunSave.bodyVars.concat(v)
      })
    }
    this.dataRunSave.bodyMode = result.data["bodyMode"]

    this.dataRunSave.headerVars = []
    if (result.data["headerVars"]) {
      _.forEach(result.data["headerVars"], v => {
        this.dataRunSave.headerVars = this.dataRunSave.headerVars.concat(v)
      })
    }

    this.dataRunSave.respVars = []
    if (result.data["respVars"]) {
      _.forEach(result.data["respVars"], v => {
        this.dataRunSave.respVars = this.dataRunSave.respVars.concat(v)
      })
    }

    this.dataRunSave.actions = []
    if (result.data["actions"]) {
      _.forEach(result.data["actions"], v => {
        this.dataRunSave.actions = this.dataRunSave.actions.concat(v)
      })
    }

    this.dataRunSave.asserts = []
    if (result.data["asserts"]) {
      _.forEach(result.data["asserts"], v => {
        this.dataRunSave.asserts = this.dataRunSave.asserts.concat(v)
      })
    }

    this.dataRunSave.preApis = []
    if (result.data["preApis"]) {
      _.forEach(result.data["preApis"], v => {
        this.dataRunSave.preApis = this.dataRunSave.preApis.concat(v)
      })
    }

    this.dataRunSave.postApis = []
    if (result.data["postApis"]) {
      _.forEach(result.data["postApis"], v => {
        this.dataRunSave.postApis = this.dataRunSave.postApis.concat(v)
      })
    }

    this.dataRunSave.otherConfigs = []
    if (result.data["otherConfig"]) {
      _.forEach(result.data["otherConfig"], v => {
        this.dataRunSave.otherConfigs = this.dataRunSave.otherConfigs.concat(v)
      })
    }
  }

  async getApiDataDetail() {
    let appName = this.apiRunSave.app
    let module = this.apiRunSave.module
    let apiDesc = this.apiRunSave.apiDesc
    let dataDesc = this.apiRunSave.dataDesc
    let method = this.apiRunSave.method
    let path = this.apiRunSave.path

    this.reqDataRespList.response = ""
    if (!dataDesc) {
      return
    }
    if (dataDesc.length > 0){
      let result
      if (appName.length>0) {
        result = await API.get<Req.ResponseModel[]>('/appList/'+appName+"?module="+module+"&apiDesc="+apiDesc+"&dataDesc="+dataDesc+"&method="+method+"&path="+path+"&mode=run")
      } else {
        result = await API.get<Req.ResponseModel[]>('/dataList'+"?dataDesc="+dataDesc)
      }
      if (result.data) {
        if (this.apiRunSave.prefix.length===0) {
          this.apiRunSave.prefix = result.data["prefix"]
        }
        if (this.apiRunSave.app.length===0) {
          this.apiRunSave.app = result.data["app"]
        }

        if (this.apiRunSave.method.length===0) {
          this.apiRunSave.method = result.data["method"]
        }

        if (this.runMethodOptions.indexOf(this.apiRunSave.method)<0) {
          this.runMethodOptions.push(this.apiRunSave.method)
        }

        if (this.apiRunSave.path.length===0) {
          this.apiRunSave.path = result.data["path"]
        }

        if (this.runApiOptions.indexOf(this.apiRunSave.path)<0) {
          this.runApiOptions.push(this.apiRunSave.path)
        }

        if (this.apiRunSave.apiDesc.length===0) {
          this.apiRunSave.apiDesc = result.data["apiDesc"]
        }

        if (this.runApiDescOptions.indexOf(this.apiRunSave.apiDesc)<0) {
          this.runApiDescOptions.push(this.apiRunSave.apiDesc)
        }

        if (this.apiRunSave.module.length===0) {
          this.apiRunSave.module = result.data["module"]
        }
        if (this.runModuleOptions.indexOf(this.apiRunSave.module)<0) {
          this.runModuleOptions.push(this.apiRunSave.module)
        }

        if (result.data["pathVars"]) {
          this.apiRunSave.pathVars = []
          _.forEach(result.data["pathVars"], v => {
            this.apiRunSave.pathVars = this.apiRunSave.pathVars.concat(v)

          })
        }

        if (result.data["queryVars"]) {
          this.apiRunSave.queryVars = []
          _.forEach(result.data["queryVars"], v => {
            this.apiRunSave.queryVars = this.apiRunSave.queryVars.concat(v)
          })
        }

        if (result.data["bodyVars"]) {
          this.apiRunSave.bodyVars = []
          _.forEach(result.data["bodyVars"], v => {
            this.apiRunSave.bodyVars = this.apiRunSave.bodyVars.concat(v)
          })
        }
        this.apiRunSave.bodyMode = result.data["bodyMode"]


        if (result.data["headerVars"]) {
          this.apiRunSave.headerVars = []
          _.forEach(result.data["headerVars"], v => {
            this.apiRunSave.headerVars = this.apiRunSave.headerVars.concat(v)
          })
        }

        if (result.data["respVars"]) {
          this.apiRunSave.respVars = []
          _.forEach(result.data["respVars"], v => {
            this.apiRunSave.respVars = this.apiRunSave.respVars.concat(v)
          })
        }


        if (result.data["actions"]) {
          this.apiRunSave.actions = []
          _.forEach(result.data["actions"], v => {
            this.apiRunSave.actions = this.apiRunSave.actions.concat(v)
          })
        }


        if (result.data["asserts"]) {
          this.apiRunSave.asserts = []
          _.forEach(result.data["asserts"], v => {
            this.apiRunSave.asserts = this.apiRunSave.asserts.concat(v)
          })
        }


        if (result.data["preApis"]) {
          this.apiRunSave.preApis = []
          _.forEach(result.data["preApis"], v => {
            this.apiRunSave.preApis = this.apiRunSave.preApis.concat(v)
          })
        }

        if (result.data["postApis"]) {
          this.apiRunSave.postApis = []
          _.forEach(result.data["postApis"], v => {
            this.apiRunSave.postApis = this.apiRunSave.postApis.concat(v)
          })
        }


        if (result.data["otherConfig"]) {
          this.apiRunSave.otherConfigs = []
          _.forEach(result.data["otherConfig"], v => {
            this.apiRunSave.otherConfigs = this.apiRunSave.otherConfigs.concat(v)
          })
        }
      }
    }
  }

  async getApiPathData() {
    this.reqDataRespList.response = ""
    this.apiRunSave.dataDesc = ""
    let appName = ""
    let path = ""
    let method = ""
    if (this.domainMode==="def") {
      if (this.apiDefSave.app) {
        appName = this.apiDefSave.app
      }
      if (this.apiDefSave.path) {
        path = this.apiDefSave.path
      }
      if (this.apiDefSave.method) {
        method = this.apiDefSave.method
      }
    } else if (this.domainMode==="run") {
      if (this.apiRunSave.app) {
        appName = this.apiRunSave.app
      }
      if (this.apiRunSave.path) {
        path = this.apiRunSave.path
      }
      if (this.apiRunSave.method) {
        method = this.apiRunSave.method
      }
    }

    if (appName.length>0 && path.length > 0 && method.length > 0){
      let result = await API.get<Req.ResponseModel[]>('/appList/'+appName+"?"+"method="+method+"&"+"path="+path+"&mode="+this.domainMode)
      if (result.data) {
        if (this.domainMode==="def") {
          this.apiDefSave.apiDesc = result.data["apiDesc"]
          this.apiDefSave.module = result.data["module"]
          // this.dataDescOptions = result.data["datasDesc"]
          this.apiDefSave.pathVars = []
          if (result.data["pathVars"]) {
            _.forEach(result.data["pathVars"], v => {
              this.apiDefSave.pathVars = this.apiDefSave.pathVars.concat(v)
            })
          }
          this.apiDefSave.queryVars = []
          if (result.data["queryVars"]) {
            _.forEach(result.data["queryVars"], v => {
              this.apiDefSave.queryVars = this.apiDefSave.queryVars.concat(v)
            })
          }
          this.apiDefSave.bodyVars = []
          if (result.data["bodyVars"]) {
            _.forEach(result.data["bodyVars"], v => {
              this.apiDefSave.bodyVars = this.apiDefSave.bodyVars.concat(v)
            })
          }
          this.apiDefSave.headerVars = []
          if (result.data["headerVars"]) {
            _.forEach(result.data["headerVars"], v => {
              this.apiDefSave.headerVars = this.apiDefSave.headerVars.concat(v)
            })
          }
          this.apiDefSave.respVars = []
          if (result.data["respVars"]) {
            _.forEach(result.data["respVars"], v => {
              this.apiDefSave.respVars = this.apiDefSave.respVars.concat(v)
            })
          }
        } else if (this.domainMode==="run") {
          this.apiRunSave.apiDesc = result.data["apiDesc"]
          this.apiRunSave.module = result.data["module"]
          this.apiRunSave.pathVars = []
          if (result.data["pathVars"]) {
            _.forEach(result.data["pathVars"], v => {
              this.apiRunSave.pathVars = this.apiRunSave.pathVars.concat(v)
            })
          }
          this.apiRunSave.queryVars = []
          if (result.data["queryVars"]) {
            _.forEach(result.data["queryVars"], v => {
              this.apiRunSave.queryVars = this.apiRunSave.queryVars.concat(v)
            })
          }
          this.apiRunSave.bodyVars = []
          if (result.data["bodyVars"]) {
            _.forEach(result.data["bodyVars"], v => {
              this.apiRunSave.bodyVars = this.apiRunSave.bodyVars.concat(v)
            })
          }
          this.apiRunSave.headerVars = []
          if (result.data["headerVars"]) {
            _.forEach(result.data["headerVars"], v => {
              this.apiRunSave.headerVars = this.apiRunSave.headerVars.concat(v)
            })
          }
          this.apiRunSave.respVars = []
          if (result.data["respVars"]) {
            _.forEach(result.data["respVars"], v => {
              this.apiRunSave.respVars = this.apiRunSave.respVars.concat(v)
            })
          }

          this.apiRunSave.otherConfigs = []
          if (result.data["otherConfig"]) {
            _.forEach(result.data["otherConfig"], v => {
              this.apiRunSave.otherConfigs = this.apiRunSave.otherConfigs.concat(v)
            })
          }
        }
      }
    }
  }

  async apiRun() {
    this.isSending = true
    this.isResquest = true
    this.isResponse = true
    let data = new URLSearchParams()
    data.append('app', JSON.stringify(this.apiRunSave.app))
    data.append('module', JSON.stringify(this.apiRunSave.module))
    data.append('apiDesc', JSON.stringify(this.apiRunSave.apiDesc))
    data.append('dataDesc', JSON.stringify(this.apiRunSave.dataDesc))
    data.append('method', JSON.stringify(this.apiRunSave.method))
    data.append('prototype', JSON.stringify(this.apiRunSave.prototype))
    data.append('path', JSON.stringify(this.apiRunSave.path))
    data.append('prefix', JSON.stringify(this.apiRunSave.prefix))
    data.append('host', JSON.stringify(this.apiRunSave.hostIp))
    data.append('pathVars', JSON.stringify(this.apiRunSave.pathVars))
    data.append('queryVars', JSON.stringify(this.apiRunSave.queryVars))
    data.append('bodyVars', JSON.stringify(this.apiRunSave.bodyVars))
    data.append('headerVars', JSON.stringify(this.apiRunSave.headerVars))
    data.append('respVars', JSON.stringify(this.apiRunSave.respVars))
    data.append('actions', JSON.stringify(this.apiRunSave.actions))
    data.append('asserts', JSON.stringify(this.apiRunSave.asserts))
    data.append('preApis', JSON.stringify(this.apiRunSave.preApis))
    data.append('postApis', JSON.stringify(this.apiRunSave.postApis))
    data.append('product', JSON.stringify(this.apiRunSave.product))
    data.append('otherConfig', JSON.stringify(this.apiRunSave.otherConfigs))
    if (this.apiRunSave.bodyMode == 'raw') {
      data.append('bodyMode', JSON.stringify(this.rawRunContentType))
    } else {
      data.append('bodyMode', JSON.stringify(this.bodyMode))
    }

    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/dataRun', data)
    if (result.data) {
      this.reqDataRespList.response = result.data["response"]
      this.reqDataRespList.url = result.data["url"]
      this.reqDataRespList.header = result.data["header"]
      this.reqDataRespList.request = result.data["request"]
      this.reqDataRespList.testResult = result.data["testResult"]
      this.reqDataRespList.failReason = result.data["failReason"]
      this.reqDataRespList.output = result.data["output"]
    } else {
      this.reqDataRespList.response = ""
      this.reqDataRespList.url = ""
      this.reqDataRespList.header = ""
      this.reqDataRespList.request = ""
      this.reqDataRespList.testResult = ""
      this.reqDataRespList.failReason = ""
      this.reqDataRespList.output = ""
    }

    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 3,
        content: result.msg
      })
    }
    this.isSending = false
  }

  async historyRun() {
    this.isSending = true
    this.isResquest = true
    this.isResponse = true
    let data = new URLSearchParams()
    data.append('app', JSON.stringify(this.historyRunSave.app))
    data.append('module', JSON.stringify(this.historyRunSave.module))
    data.append('apiDesc', JSON.stringify(this.historyRunSave.apiDesc))
    data.append('dataDesc', JSON.stringify(this.historyRunSave.dataDesc))
    data.append('prototype', JSON.stringify(this.historyRunSave.prototype))
    data.append('method', JSON.stringify(this.historyRunSave.method))
    data.append('path', JSON.stringify(this.historyRunSave.path))
    data.append('host', JSON.stringify(this.historyRunSave.host))
    data.append('prefix', JSON.stringify(this.historyRunSave.prefix))
    data.append('pathVars', JSON.stringify(this.historyRunSave.pathVars))
    data.append('queryVars', JSON.stringify(this.historyRunSave.queryVars))
    data.append('bodyVars', JSON.stringify(this.historyRunSave.bodyVars))
    data.append('headerVars', JSON.stringify(this.historyRunSave.headerVars))
    data.append('respVars', JSON.stringify(this.historyRunSave.respVars))
    data.append('actions', JSON.stringify(this.historyRunSave.actions))
    data.append('asserts', JSON.stringify(this.historyRunSave.asserts))
    data.append('preApis', JSON.stringify(this.historyRunSave.preApis))
    data.append('postApis', JSON.stringify(this.historyRunSave.postApis))
    data.append('product', JSON.stringify(this.historyRunSave.product))
    data.append('fileName', JSON.stringify(this.historyRunSave.fileName))
    data.append('otherConfig', JSON.stringify(this.historyRunSave.otherConfigs))
    data.append('bodyMode', JSON.stringify(this.historyRunSave.bodyMode))

    // console.log("this.historyRunSave.bodyMode: ", this.historyRunSave.bodyMode)
    // if (this.historyRunSave.bodyMode == 'raw') {
    //   data.append('bodyMode', JSON.stringify(this.rawHistoryContentType))
    // } else {
    //   data.append('bodyMode', JSON.stringify(this.bodyMode))
    // }
    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/historyRun', data)
    if (result.data) {
      this.historyModeReqDataRespList.response = result.data["response"]
      this.historyModeReqDataRespList.url = result.data["url"]
      this.historyModeReqDataRespList.header = result.data["header"]
      this.historyModeReqDataRespList.request = result.data["request"]
      this.historyModeReqDataRespList.testResult = result.data["testResult"]
      this.historyModeReqDataRespList.failReason = result.data["failReason"]
    } else {
      this.historyModeReqDataRespList.response = ""
      this.historyModeReqDataRespList.url = ""
      this.historyModeReqDataRespList.header = ""
      this.historyModeReqDataRespList.request = ""
      this.historyModeReqDataRespList.testResult = ""
      this.historyModeReqDataRespList.failReason = ""
    }

    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 3,
        content: result.msg
      })
    }
    this.isSending = false
    this.getHistoryMenu()
  }

  async dataRun() {
    this.isSending = true
    this.isResquest = true
    this.isResponse = true
    let data = new URLSearchParams()
    data.append('app', JSON.stringify(this.dataRunSave.app))
    data.append('module', JSON.stringify(this.dataRunSave.module))
    data.append('apiDesc', JSON.stringify(this.dataRunSave.apiDesc))
    data.append('dataDesc', JSON.stringify(this.dataRunSave.dataDesc))
    data.append('method', JSON.stringify(this.dataRunSave.method))
    data.append('prototype', JSON.stringify(this.dataRunSave.prototype))
    data.append('path', JSON.stringify(this.dataRunSave.path))
    data.append('prefix', JSON.stringify(this.dataRunSave.prefix))
    data.append('host', JSON.stringify(this.dataRunSave.hostIp))
    data.append('pathVars', JSON.stringify(this.dataRunSave.pathVars))
    data.append('queryVars', JSON.stringify(this.dataRunSave.queryVars))
    data.append('bodyVars', JSON.stringify(this.dataRunSave.bodyVars))
    data.append('headerVars', JSON.stringify(this.dataRunSave.headerVars))
    data.append('respVars', JSON.stringify(this.dataRunSave.respVars))
    data.append('actions', JSON.stringify(this.dataRunSave.actions))
    data.append('asserts', JSON.stringify(this.dataRunSave.asserts))
    data.append('preApis', JSON.stringify(this.dataRunSave.preApis))
    data.append('postApis', JSON.stringify(this.dataRunSave.postApis))
    data.append('product', JSON.stringify(this.dataRunSave.product))
    data.append('otherConfig', JSON.stringify(this.dataRunSave.otherConfigs))
    data.append('bodyMode', JSON.stringify(this.dataRunSave.bodyMode))
    // if (this.dataRunSave.bodyMode == 'raw') {
    //   data.append('bodyMode', JSON.stringify(this.rawDataContentType))
    // } else {
    //   data.append('bodyMode', JSON.stringify(this.bodyMode))
    // }
    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/dataRun', data)
    if (result.data) {
      this.dataModeReqDataRespList.response = result.data["response"]
      this.dataModeReqDataRespList.url = result.data["url"]
      this.dataModeReqDataRespList.header = result.data["header"]
      this.dataModeReqDataRespList.request = result.data["request"]
      this.dataModeReqDataRespList.testResult = result.data["testResult"]
      this.dataModeReqDataRespList.failReason = result.data["failReason"]
      this.dataModeReqDataRespList.output = result.data["output"]
    } else {
      this.dataModeReqDataRespList.response = ""
      this.dataModeReqDataRespList.url = ""
      this.dataModeReqDataRespList.header = ""
      this.dataModeReqDataRespList.request = ""
      this.dataModeReqDataRespList.testResult = ""
      this.dataModeReqDataRespList.failReason = ""
      this.dataModeReqDataRespList.output = ""
    }

    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 3,
        content: result.msg
      })
    }
    this.isSending = false
  }

  async sceneRun() {
    this.isSending = true
    this.isResquest = true
    this.isResponse = true
    let data = new URLSearchParams()
    data.append('product', JSON.stringify(this.sceneSave.product))
    data.append('name', JSON.stringify(this.sceneSave.name))
    data.append('dataList', JSON.stringify(this.sceneSave.dataList))
    data.append('type', JSON.stringify(this.sceneSave.type))
    data.append('runNum', JSON.stringify(this.sceneSave.runNum))

    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/sceneRun', data)
    if (result.data) {
      this.sceneModeReqDataRespList.lastDataFile = result.data["lastDataFile"]
      this.sceneModeReqDataRespList.testResult = result.data["testResult"]
      this.sceneModeReqDataRespList.failReason = result.data["failReason"]
    } else {
      this.sceneModeReqDataRespList.lastDataFile = ""
      this.sceneModeReqDataRespList.testResult = ""
      this.sceneModeReqDataRespList.failReason = ""
    }

    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 3,
        content: result.msg
      })
    }
    this.isSending = false
  }

  async sceneHistoryRun() {
    this.isSending = true
    this.isResquest = true
    this.isResponse = true
    let data = new URLSearchParams()
    data.append('product', JSON.stringify(this.sceneHistorySave.product))
    data.append('name', JSON.stringify(this.sceneHistorySave.name))
    data.append('dataList', JSON.stringify(this.sceneHistorySave.dataList))
    data.append('type', JSON.stringify(this.sceneHistorySave.type))
    data.append('runNum', JSON.stringify(this.sceneHistorySave.runNum))

    this.start = moment().valueOf()
    let result = await API.post<Req.ResponseModel[]>('/sceneRun', data)
    if (result.data) {
      this.sceneHistoryRespList.lastDataFile = result.data["lastDataFile"]
      this.sceneHistoryRespList.testResult = result.data["testResult"]
      this.sceneHistoryRespList.failReason = result.data["failReason"]
    } else {
      this.sceneHistoryRespList.lastDataFile = ""
      this.sceneHistoryRespList.testResult = ""
      this.sceneHistoryRespList.failReason = ""
    }

    if (result.code == 200) {
      this.$Message.success({
        duration: 3,
        content: result.msg
      })
    } else {
      this.$Message.error({
        duration: 3,
        content: result.msg
      })
    }
    this.isSending = false
  }

  onAddApi(val) {
    if (this.domainMode==="def") {
      this.defApiOptions.push(val)
    } else  if (this.domainMode==="run") {
      this.runApiOptions.push(val)
    }
  }

  onAddMethod(val) {
    if (this.domainMode==="def") {
      this.defMethodOptions.push(val)
    } else  if (this.domainMode==="run") {
      this.runMethodOptions.push(val)
    }
  }

  onAddApp(val) {
    this.appOptions.push(val)
  }

  onAddPlaybook(val) {
    this.playbookOptions.push(val)
  }

  onAddModule(val) {
    if (this.domainMode==="def") {
      this.defModuleOptions.push(val)
    } else  if (this.domainMode==="run") {
      this.runModuleOptions.push(val)
    }

  }

  onAddApiDesc(val) {
    if (this.domainMode==="def") {
      this.defApiDescOptions.push(val)
    } else  if (this.domainMode==="run") {
      this.runApiDescOptions.push(val)
    }
  }

  onAddDataDesc(val) {
    if (!this.apiRunSave.dataDesc){
      this.apiRunSave.dataDesc = val
    }
    this.dataDescOptions.push(val)
  }

  async getEnvList() {
    let result = await API.get<Req.ResponseModel[]>('/envList')
    if (result.data) {
      this.depDataTableData = result.data.map(item => (item['product']));
    }
  }

  async getApiEnv(val) {
    let result = await API.get<Req.ResponseModel[]>('/envList/'+val)
    if (result.data) {
      this.apiRunSave.prototype = result.data["protocol"];
      this.apiRunSave.hostIp = result.data["ip"];
      this.apiRunSave.headerVars = result.data["auth"];
    }
  }

  async getHistoryEnv(val) {
    let result = await API.get<Req.ResponseModel[]>('/envList/'+val)
    if (result.data) {
      this.historyRunSave.prototype = result.data["protocol"];
      this.historyRunSave.host = result.data["ip"];
      this.historyRunSave.headerVars = result.data["auth"];
    }
  }

  async getDataEnv(val) {
    let result = await API.get<Req.ResponseModel[]>('/envList/'+val)
    if (result.data) {
      this.dataRunSave.prototype = result.data["protocol"];
      this.dataRunSave.hostIp = result.data["ip"];
      this.dataRunSave.headerVars = result.data["auth"];
    }
  }

  async getSceneEnv(val) {
    let result = await API.get<Req.ResponseModel[]>('/envList/'+val)
    if (result.data) {
      this.sceneSave.product = result.data["product"];
    }
  }

  async getDataFileList() {
    let result = await API.get<Req.ResponseModel[]>('/dataFileList')
    if (result.data) {
      this.allDataFile = result.data.map(item => (item['dataFile']));
    }
  }
}
</script>
