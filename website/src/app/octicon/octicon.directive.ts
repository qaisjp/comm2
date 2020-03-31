import { Directive, Input, OnInit, ElementRef, Renderer2 } from '@angular/core';

import * as octicons from '@primer/octicons';

@Directive({
  selector: '[appOcticon]',
})
export class OcticonDirective implements OnInit {
  @Input() appOcticon: string;
  @Input() color: string;
  @Input() width: number;

  constructor(private elementRef: ElementRef, private renderer: Renderer2) { }

  ngOnInit(): void {
    const el: HTMLElement = this.elementRef.nativeElement;
    el.innerHTML = octicons[this.appOcticon].toSVG();
    el.classList.add('octicon-parent');

    // tslint:disable-next-line:no-non-null-assertion
    const icon: ChildNode = el.firstChild!;
    if (this.color) {
      this.renderer.setStyle(icon, 'fill', this.color);
    }
    if (this.width) {
      this.renderer.setStyle(icon, 'width', this.width);
      this.renderer.setStyle(icon, 'height', '100%');
    }
  }
}
