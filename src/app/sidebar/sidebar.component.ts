import { Component } from '@angular/core';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.css']
})

export class SidebarComponent {
  sideBarItems: SideBarItem[] = [
    { icon: 'bi-hdd', text: 'My Drive', routerLink: '/folderFile' },
    { icon: 'bi-search', text: 'Search', routerLink: '/search' },
    { icon: 'bi-person-gear', text: 'Admin Panel', routerLink: '' },
    { icon: 'bi-gear', text: 'Settings', routerLink: '' },
    { icon: 'bi-activity', text: 'System Monitor', routerLink: '' },
  ];

  selectedItem: SideBarItem | null = this.sideBarItems[0];
  singleClick(item: SideBarItem) {
    this.selectedItem = item;
  }
}

interface SideBarItem {
  icon: string;
  text: string;
  routerLink: string;
}
