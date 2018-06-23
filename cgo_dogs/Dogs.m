#import "Dogs.h"

@interface Dog : RLMObject
@property NSString *name;
@property NSData   *picture;
@property NSInteger age;
@end
@implementation Dog
@end
RLM_ARRAY_TYPE(Dog)
@interface Person : RLMObject
@property NSString             *name;
@property RLMArray<Dog *><Dog> *dogs;
@end
@implementation Person
@end


const NSString* CountDogs() {
    Dog *mydog = [[Dog alloc] init];
    mydog.name = @"Rex";
    mydog.age = 1;
    mydog.picture = nil; // properties are nullable
    NSLog(@"Name of dog: %@", mydog.name);
    // Query Realm for all dogs less than 2 years old
    RLMResults<Dog *> *puppies = [Dog objectsWhere:@"age < 2"];
    puppies.count; // => 0 because no dogs have been added to the Realm yet
    // Persist your data easily
    RLMRealm *realm = [RLMRealm defaultRealm];
        [realm transactionWithBlock:^{
        [realm addObject:mydog];
    }];
    // Queries are updated in realtime
    puppies.count; // => 1
    NSLog(@"Dogs Count: %d", puppies.count);
    return @"";
} // CountDogs