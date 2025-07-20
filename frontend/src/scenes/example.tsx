import {
  makeScene2D,
  SVG,
  Path
} from '@motion-canvas/2d';
import {
  all,
  waitFor,
  createRef
} from '@motion-canvas/core';
import dog from '/public/dog_1.svg?raw';

export default makeScene2D(function* (view) {
  view.fill('#000000');

  const dogRef = createRef<SVG>();

  view.add(
    <SVG 
      ref={dogRef} 
      svg={dog.replace("@color", "#f2ff48")} 
      size={300} 
      position={() => [view.width() / 2, view.height() / 2]}
    />
  );

  yield* dogRef().scale(1, 0.5);  
  yield* dogRef().opacity(1, 0.5); 

  for (const child of dogRef().children()[0].children()) {
    if (child instanceof Path) {
      yield* child.fill('white', 1);
    }
  }
});
