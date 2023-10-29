#import <Foundation/Foundation.h>

#import "../../../ddns4cdn.h"

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        NSUserDefaults *args = [NSUserDefaults standardUserDefaults];
        NSString *config = [args stringForKey:@"c"] ?: @"config.toml";

        NSError *error = nil;
        NSString *data = [NSString stringWithContentsOfFile:config encoding:NSUTF8StringEncoding error:&error];
        if (error) {
            NSLog(@"%@", error.localizedDescription);
            return 1;
        }

        Worker((char *)[data UTF8String]);
    }
    return 0;
}
