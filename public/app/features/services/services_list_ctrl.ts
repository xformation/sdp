import angular from 'angular';
import { BASE_URL } from 'app/main.conf';

export class ServicesCtrl {
  activeServices: any[];
  activeTabIndex: number;
  navModel: any;
  tabs: any;
  serviceList: any[];
  searchText: string;
  constructor($location, private $http: any) {
    // this.navModel = navModelSrv.getServicesNav();
    this.searchText = '';
    this.serviceList = [];
    this.getServiceList().then(resData => {
      this.serviceList = resData.data;
      this.tabs = this.getTabsFromData(this.serviceList);
      this.changeTab(this.tabs[0]);
    });
  }

  getServiceList() {
    return this.$http.get(BASE_URL + 'services');
  }

  getTabsFromData(data) {
    let tabs = [];
    let retData = [];
    for (var i in data) {
      let service = data[i];
      let tabIndex = tabs.indexOf(service.category['name']);
      if (tabIndex === -1) {
        tabIndex = tabs.length;
        tabs.push(service.category['name']);
        retData.push({
          label: service.category['name'],
          tabIndex: tabIndex,
        });
      }
      service['tabIndex'] = tabIndex;
    }
    return retData;
  }

  changeTab(tab) {
    this.activeTabIndex = tab.tabIndex;
    this.filterServices(this.activeTabIndex, this.searchText);
  }

  filterServices(activeTabIndex, searchText) {
    this.activeServices = [];
    this.activeServices = this.serviceList.filter(service => {
      let regex = new RegExp(searchText, 'ig');
      return service.tabIndex === activeTabIndex && regex.test(service.name);
    });
  }

  onSearchService() {
    this.filterServices(this.activeTabIndex, this.searchText);
  }

  filterServicesOnText(name) {
    this.activeServices = [];
    this.activeServices = this.serviceList.filter(function(service) {
      let regex = new RegExp(name, 'ig');
      return regex.test(service.name);
    });
  }

  subscribe(service) {
    let subscriptionData = {
      id: '1',
      serviceId: service.id,
      customerId: '12',
    };
    this.$http.post('http://10.10.10.60:8080/api/v1/subscription/', subscriptionData).then(data => {
      alert('Service subscribed successfully');
    });
  }
}
angular.module('grafana.controllers').controller('ServicesCtrl', ServicesCtrl);
