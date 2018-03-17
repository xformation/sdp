import angular from 'angular';
import { BASE_URL } from 'app/main.conf';
export class HomeCtrl {
  navModel: any;
  resourceData: any;
  recentAssets: any;
  availableResources: any;
  topPerformingAssets: any;
  projectCentralData: any;
  lineChartData: any;
  constructor($scope, private $http: any) {
    this.resourceData = [];
    this.recentAssets = [];
    this.availableResources = [];
    this.topPerformingAssets = [];
    // this.navModel = navModelSrv.getHomeNav();
    this.getResourceCentralData().then(res => {
      this.resourceData = res.data;
    });
    this.getRecentAssets().then(res => {
      this.recentAssets = res.data;
    });
    this.getAvailableResources().then(res => {
      this.availableResources = res.data;
    });
    this.getTopPerformingAssets().then(res => {
      this.topPerformingAssets = res.data;
    });
    // this.getProjectCentralData().then(res => {
    //     this.projectCentralData = JSON.stringify(res.data);
    // });
    this.projectCentralData = JSON.stringify(this.getProjectCentralData());
    this.lineChartData = JSON.stringify(this.getLineChartData());
  }

  getTopPerformingAssets() {
    return this.$http.get(BASE_URL + 'top_performing_resources');
  }

  getResourceCentralData() {
    return this.$http.get(BASE_URL + 'resources');
  }

  getRecentAssets() {
    return this.$http.get(BASE_URL + 'recent_assets');
  }

  getAvailableResources() {
    return this.$http.get(BASE_URL + 'available_resources');
  }

  getProjectCentralData() {
    // return this.$http.get(BASE_URL + "project_central_data");
    return [
      {
        label: 'P1',
        data: [20],
        backgroundColor: ['#0ecbf5'],
      },
      {
        label: 'P2',
        data: [15],
        backgroundColor: ['#028bff'],
      },
      {
        label: 'P3',
        data: [35],
        backgroundColor: ['#0064ba'],
      },
    ];
  }

  getLineChartData() {
    return [
      {
        label: [1, 2, 3, 4],
        data: [20, 40, 50, 30],
        fill: false,
        backgroundColor: 'rgb(0,0,0)',
        pointBackgroundColor: '#0064ba',
        borderColor: '#0064ba',
      },
    ];
  }
}

angular.module('grafana.controllers').controller('HomeCtrl', HomeCtrl);
