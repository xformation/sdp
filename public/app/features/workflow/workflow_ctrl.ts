import angular from 'angular';
import { BASE_URL } from 'app/main.conf';
export class WorkflowCtrl {
  navModel: any;
  workflowData: any;
  activeTabIndex: number;
  activeStepIndex: number;
  activeSubstepIndex: number;
  tabs: any;
  activeStepData: any;
  activeTabData: any;
  activeSubstepData: any;
  substeps: any;
  isStepsAvailable: boolean;
  constructor(private $http: any, $location) {
    // this.navModel = navModelSrv.getWorkflowNav();
    this.tabs = [];
    this.substeps = [];
    this.isStepsAvailable = false;
    this.getWorkflowData().then(resData => {
      this.workflowData = resData.data;
      this.tabs = this.getTabs();
      this.changeTab(0);
    });
  }

  changeTab(tabIndex) {
    this.activeTabIndex = tabIndex;
    this.setActiveTabData(tabIndex);
  }

  changeStep(stepIndex) {
    if (this.activeTabData.steps) {
      this.activeStepIndex = stepIndex;
      this.setActiveStepData(stepIndex);
      this.isStepsAvailable = true;
    } else {
      this.isStepsAvailable = false;
    }
  }

  changeSubstep(stepIndex) {
    this.activeSubstepIndex = stepIndex;
    this.setActiveSubStepData(stepIndex);
  }

  setActiveSubStepData(stepIndex) {
    this.activeSubstepData = this.substeps[stepIndex];
  }

  setActiveTabData(tabIndex) {
    this.activeTabData = this.workflowData[tabIndex];
    this.changeStep(0);
  }

  setActiveStepData(stepInex) {
    this.activeStepData = this.activeTabData.steps[stepInex];
    this.substeps = this.activeStepData.substeps;
    this.changeSubstep(0);
  }

  getWorkflowData() {
    return this.$http.get(BASE_URL + 'workflow');
  }

  getTabs() {
    let tabs = [];
    for (var i in this.workflowData) {
      let info = this.workflowData[i];
      tabs.push({
        label: info.tab,
        class: info.class,
        tabIndex: i,
      });
    }
    return tabs;
  }
}
angular.module('grafana.controllers').controller('WorkflowCtrl', WorkflowCtrl);
